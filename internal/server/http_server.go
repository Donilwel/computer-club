package server

import (
	"computer-club/internal/di"
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
