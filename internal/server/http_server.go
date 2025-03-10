package server

import (
	"computer-club/internal/di"
	"context"
	"fmt"
	"net/http"
)

// HttpServer управляет HTTP-сервером
type HttpServer struct {
	container *di.Container
	server    *http.Server
}

// NewHttpServer создаёт HTTP-сервер
func NewHttpServer(container *di.Container) *HttpServer {
	httpSrv := &http.Server{
		Addr:    ":" + container.Cfg.Server.HTTPPort,
		Handler: container.Router,
	}

	return &HttpServer{
		container: container,
		server:    httpSrv,
	}
}

// Run запускает HTTP-сервер
func (s *HttpServer) Run() {
	fmt.Println("HTTP сервер запущен на порту:", s.container.Cfg.Server.HTTPPort)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.container.Log.Fatal("Ошибка HTTP сервера:", err)
	}
}

// Shutdown корректно завершает работу HTTP-сервера
func (s *HttpServer) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
	s.container.Log.Info("HTTP сервер успешно остановлен")
}
