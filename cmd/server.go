package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"short-link/api/app"
	"short-link/api/middleware"
	"short-link/config"
	"short-link/logs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

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
