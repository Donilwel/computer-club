package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"computer-club/config"
	"computer-club/internal/repository"
	"computer-club/internal/service"
	"computer-club/internal/transport"
	pb "computer-club/proto/auth"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	repo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(repo, cfg.JWT.Secret, time.Hour*time.Duration(cfg.JWT.ExpirationHours))
	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, transport.NewAuthServer(authService))

	listener, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	log.Println("gRPC сервер запущен на порту", cfg.Server.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC: %v", err)
	}
}
