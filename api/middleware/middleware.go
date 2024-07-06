package middleware

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"short-link/internal/metrics"
	"short-link/logs"
	"short-link/utils/apix"
	"short-link/utils/gox"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var publicEndPoint = map[string]bool{
	"/api/v1/admin/users/register": true,
	"/api/v1/admin/users/login":    true,
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logs.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("IP", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPIPe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPIPe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPIPe {
					logs.Error(nil, c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) //nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logs.Error(nil, "[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logs.Error(nil, "[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		if isPublicEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := apix.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "userId", mc["userId"]))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "username", mc["username"]))
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

func IP() gin.HandlerFunc {
	return func(c *gin.Context) {
		IP := c.ClientIP()
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "IP", IP))
		c.Next()
	}
}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			shortUrl := c.Param("short-link")
			if c.Writer.Status() == http.StatusFound && shortUrl != "" {
				// 成功转发
				sr := metrics.ShortUrlRequest{
					ShortUrl: shortUrl,
					IP:       c.ClientIP(),
				}
				gox.Run(context.Background(), func(ctx context.Context) {
					metrics.RecordShortUrlRequest(&sr)
				})
			}
		}()
		c.Next()
	}
}

func isPublicEndpoint(path string) bool {
	return publicEndPoint[path]
}
