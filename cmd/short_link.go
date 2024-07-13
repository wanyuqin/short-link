package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"short-link/config"
	"short-link/database/cache"
	"short-link/database/mysql"
	"short-link/logs"
	"time"

	"github.com/spf13/cobra"
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
		Run: func(cmd *cobra.Command, args []string) {
			startEventBus()

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
