package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"short-link/api"
	"short-link/config"
	"time"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use: "slink",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config.InitializeConfig(cfgFile)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			startHttpServer()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "set config")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startHttpServer() {
	cfg := config.GetConfig()
	engine := gin.Default()

	gin.SetMode(cfg.Application.Mode)
	engine.Use(exceptionMiddleware)

	rootGroup := engine.Group(cfg.Application.ContextPath)
	api.NewRouter(rootGroup)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Application.Host, cfg.Application.Port),
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func exceptionMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器发生内部错误"})
		}
	}()
	c.Next()
}
