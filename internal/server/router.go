package server

import (
	"travel-agent/internal/engine"
	"travel-agent/internal/server/handler"
	"travel-agent/internal/server/middleware"
	"travel-agent/internal/session"

	"github.com/go-chi/chi/v5"
)

// NewRouter 创建路由
func NewRouter(eng engine.Engine, store session.Store, corsOrigins []string) *chi.Mux {
	r := chi.NewRouter()

	// 中间件
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(corsOrigins))

	// 健康检查
	healthHandler := handler.NewHealthHandler()
	r.Get("/health", healthHandler.Handle)

	// API 路由
	r.Route("/api", func(r chi.Router) {
		// 会话管理
		sessionHandler := handler.NewSessionHandler(store)
		r.Post("/session", sessionHandler.Create)
		r.Get("/session/{id}", sessionHandler.Get)
		r.Delete("/session/{id}", sessionHandler.Delete)

		// 聊天
		chatHandler := handler.NewChatHandler(eng, store)
		r.Post("/chat", chatHandler.Handle)
	})

	return r
}
