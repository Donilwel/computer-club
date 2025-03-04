package main

import (
  "computer-club/internal/auth"
  "computer-club/internal/config"
  "computer-club/internal/logger"
  "google.golang.org/grpc"
  "net"
)

func main() {
	cfg, _ := config.LoadConfig()
	log := logger.InitLogger(cfg.Log.Level, cfg.Log.File)

	db, _ := database.ConnectDB(cfg)
	_ = db.AutoMigrate(&auth.User{})

	authService := auth.NewAuthService(db, cfg.JWT.Secret, log)

	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authService)

	listener, _ := net.Listen("tcp", ":50051")
	log.Info("gRPC сервер запущен на порту 50051")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
