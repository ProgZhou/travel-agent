package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"travel-agent/internal/config"
	"travel-agent/internal/engine"
	"travel-agent/internal/llm"
	"travel-agent/internal/pkg"
	"travel-agent/internal/provider"
	"travel-agent/internal/server"
	"travel-agent/internal/session"
	"travel-agent/internal/tool"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	// 初始化日志
	pkg.InitLogger(cfg.Log.Level, cfg.Log.Format)

	// 检查 API Key
	if cfg.LLM.APIKey == "" {
		log.Fatal("ARK_API_KEY environment variable is required")
	}

	// 初始化组件
	llmClient := llm.NewDoubaoClient(
		cfg.LLM.APIURL,
		cfg.LLM.APIKey,
		cfg.LLM.Model,
		cfg.LLM.MaxTokens,
		cfg.LLM.Temperature,
		cfg.LLM.Timeout,
		cfg.LLM.MaxRetries,
		cfg.LLM.RetryInterval,
	)

	sessionStore := session.NewMemoryStore(cfg.Session.MaxSessions, cfg.Session.TTL)

	// 初始化 Provider
	flightProvider := provider.NewMockFlightProvider()
	hotelProvider := provider.NewMockHotelProvider()
	weatherProvider := provider.NewMockWeatherProvider()

	// 初始化 Tool Registry
	toolRegistry := tool.NewRegistry()
	toolRegistry.Register(tool.NewFlightTool(flightProvider, cfg.Tools.Flight.Timeout, cfg.Tools.Flight.RetryCount, cfg.Tools.Flight.RetryInterval))
	toolRegistry.Register(tool.NewHotelTool(hotelProvider, cfg.Tools.Hotel.Timeout, cfg.Tools.Hotel.RetryCount, cfg.Tools.Hotel.RetryInterval))
	toolRegistry.Register(tool.NewWeatherTool(weatherProvider, cfg.Tools.Weather.Timeout, cfg.Tools.Weather.RetryCount, cfg.Tools.Weather.RetryInterval))

	// 初始化 Engine
	eng := engine.NewEngine(llmClient, toolRegistry, sessionStore, cfg.Engine.MaxReactIterations)

	// 创建路由
	router := server.NewRouter(eng, sessionStore, cfg.Server.CORSOrigins)

	// 创建 HTTP 服务器
	srv := server.New(cfg.Server.Host, cfg.Server.Port, router)

	// 启动服务器
	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "error", err)
	}

	slog.Info("server stopped")
}
