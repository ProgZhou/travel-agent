package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Server HTTP 服务器
type Server struct {
	httpServer *http.Server
}

// New 创建 HTTP 服务器
func New(host string, port int, handler http.Handler) *Server {
	addr := fmt.Sprintf("%s:%d", host, port)
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	slog.Info("starting http server", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down http server")
	return s.httpServer.Shutdown(ctx)
}
