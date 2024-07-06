package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"short-link/api/admin"
	"short-link/api/app"
	"short-link/api/middleware"
	"short-link/config"
	"short-link/database/cache"
	"short-link/database/mysql"
	"short-link/docs"
	_ "short-link/docs"
	"short-link/internal/consts"
	"short-link/internal/link/event"
	"short-link/logs"
	"short-link/pkg/bus"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use: "slink",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := config.InitializeConfig(cfgFile); err != nil {
				return err
			}
			mysql.InitializeDBClient()
			cache.InitializeRedisClient()
			logs.InitializeLogger()
			return nil
		},
		// TODO 调整启动
		Run: func(cmd *cobra.Command, args []string) {
			newEvent()
			adminServer := startAdminHttpServer()
			appServer := startHttpServer()
			metricsServer := startMetricsServer()
			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
			<-quit
			logs.Info("Shutdown App Server ...")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := adminServer.Shutdown(ctx); err != nil {
				logs.Fatal("App Server Shutdown:", zap.Any("err", err))
			}
			logs.Info("Server exiting")
			if err := appServer.Shutdown(ctx); err != nil {
				logs.Fatal("App Server Shutdown:", zap.Any("err", err))
			}
			if err := metricsServer.Shutdown(ctx); err != nil {
				logs.Fatal("App Server Shutdown:", zap.Any("err", err))
			}
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "set config")
}

func Start() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// admin 服务
func startAdminHttpServer() *http.Server {
	cfg := config.GetConfig().GetHTTPConfig("admin")
	engine := gin.Default()

	gin.SetMode(cfg.Mode)
	engine.Use(middleware.GinLogger(), middleware.GinRecovery(true), middleware.CORS(), middleware.JWTAuthMiddleware())

	rootGroup := engine.Group(cfg.ContextPath)
	admin.NewRouter(rootGroup)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	docs.SwaggerInfo.BasePath = cfg.ContextPath
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	go func() {
		logs.Info("admin http server start", zap.Any("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.Fatal("listen: %s\n", zap.Any("err", err))
		}
	}()
	return srv
}

// http 服务
func startHttpServer() *http.Server {
	cfg := config.GetConfig().GetHTTPConfig("app")
	engine := gin.Default()

	gin.SetMode(cfg.Mode)
	engine.Use(middleware.IP(), middleware.GinLogger(), middleware.Metrics(), middleware.GinRecovery(true))

	rootGroup := engine.Group(cfg.ContextPath)
	app.NewRouter(rootGroup)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func() {
		logs.Info("app http server start", zap.Any("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.Fatal("listen: %s\n", zap.Any("err", err))
		}
	}()
	return srv

}

// metrics 服务
func startMetricsServer() *http.Server {
	cfg := config.GetConfig().GetHTTPConfig("metrics")
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		logs.Info("metrics http server start", zap.Any("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.Fatal("listen: %s\n", zap.Any("err", err))
		}
	}()
	return srv

}

func newEvent() {
	eventBus := bus.NewAsyncEventBus()
	eventBus.AddEventListener(consts.DeleteShortURLEvent, event.DeleteShortUrl)
}
