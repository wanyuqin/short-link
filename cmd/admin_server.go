package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"short-link/api/admin"
	"short-link/api/middleware"
	"short-link/config"
	"short-link/docs"
	"short-link/internal/consts"
	"short-link/internal/link/event"
	"short-link/logs"
	"short-link/pkg/bus"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

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

func startEventBus() {
	eventBus := bus.NewAsyncEventBus()

	eventBus.AddEventListener(consts.DeleteShortURLEvent, event.DeleteShortUrlEvent)
	eventBus.AddEventListener(consts.UpdateBlackListStatusEvent, event.UpdateBlackListStatusEvent)
	eventBus.AddEventListener(consts.DeleteBlackListEvent, event.DeleteBlackListEvent)
}
