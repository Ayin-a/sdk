package main

import (
	"context"
	"hk4e_sdk/pkg/config"
	"hk4e_sdk/pkg/http"
	"hk4e_sdk/pkg/logger"
	zq "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// 改index死妈
func main() {

	cfg := config.LoadConfig()
	logger.InitLogger(strings.ToUpper(cfg.LogLevel))

	httpsrv := http.NewServer(&cfg)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpsrv.Start(); err != nil {
			logger.Error("无法启动HTTP服务器")
		}
	}()
	if cfg.EnablePprof {
		// 原生pprof /debug/pprof
		// 可视化图表 /debug/statsviz
		go func() {
			logger.Info("性能分析正在启动，默认6060端口")
			// 将 statsviz 注册到默认的 HTTP 处理器中。
			//if err := statsviz.RegisterDefault(); err != nil {
			//	logger.Error("无法注册 statsviz 处理器")
			//}

			if err := zq.ListenAndServe("0.0.0.0:6060", nil); err != nil {
				logger.Error("无法启动性能分析")
			}
		}()
	}

	restartTicker := time.NewTicker(time.Duration(cfg.ReloadTime) * time.Second)
	go func() {
		for {
			select {
			case <-restartTicker.C:
				logger.Info("正在重启服务器...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := httpsrv.Shutdown(ctx); err != nil {
					logger.Error("无法正常关闭HTTP服务器")
				}
				// 在这里做任何需要的清理工作
				err := http.Restart()
				if err != nil {
					logger.Error("无法重启服务器")
				}
			case <-done:
				// 添加停止服务
				restartTicker.Stop()
				logger.Info("HTTP服务正在关闭")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := httpsrv.Shutdown(ctx); err != nil {
					logger.Error("无法正常关闭HTTP服务")
				}
				logger.Info("HTTP服务已停止")
				os.Exit(0) // 将终止程序

			}
		}
	}()

	// 保持main函数运行
	select {}
}
