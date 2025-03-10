package server

import (
	"computer-club/internal/di"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server управляет HTTP-сервером (и в будущем gRPC)
type Server struct {
	httpServer *http.Server
	container  *di.Container
}

// NewServer создаёт сервер с HTTP API
func NewServer() *Server {
	container := di.NewContainer()

	httpSrv := &http.Server{
		Addr:    ":" + container.Cfg.Server.HTTPPort,
		Handler: container.Router,
	}

	return &Server{
		httpServer: httpSrv,
		container:  container,
	}
}

// Run запускает HTTP-сервер и обрабатывает graceful shutdown
func (s *Server) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("HTTP сервер запущен на порту:", s.container.Cfg.Server.HTTPPort)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.container.Log.Fatal("Ошибка HTTP сервера:", err)
		}
	}()

	<-ctx.Done()

	s.GracefulShutdown()
}

// GracefulShutdown корректно завершает работу сервера
func (s *Server) GracefulShutdown() {
	s.container.Log.Info("Остановка сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.container.Log.Fatal("Ошибка при остановке HTTP сервера:", err)
	}

	s.container.Log.Info("Сервер успешно остановлен")
}
