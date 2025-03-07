package main

import (
	"computer-club/internal/logger"
	"computer-club/internal/middleware"
	"computer-club/internal/usecase"
	"context"
	"fmt"
	"net/http"

	"computer-club/config"
	"computer-club/internal/delivery/httpService"
	"computer-club/internal/repository"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger()
	// Подключаемся к БД и Redis
	db := repository.NewPostgresDB(cfg)
	redisClient := repository.NewRedisClient(cfg)
	repository.Migrate(db)

	// Создаем репозитории
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userRepo := repository.NewPostgresUserRepo(db)
	sessionRepo := repository.NewPostgresSessionRepo(db, redisClient)
	computerRepo := repository.NewComputerRepository(db)
	tariffRepo := repository.NewTariffRepositoryPostgres(db)
	walletRepo := repository.NewPostgresWalletRepo(db)

	tariffUsecase := usecase.NewTariffUsecase(tariffRepo)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, tariffUsecase)
	userUsecase := usecase.NewUserUsecase(userRepo, walletUsecase)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo, computerRepo, walletUsecase)
	computerUsecase := usecase.NewComputerUsecase(computerRepo)

	go sessionUsecase.MonitorSessions(ctx)
	// Запускаем HTTP сервер
	handler := httpService.NewHandler(userUsecase, computerUsecase, sessionUsecase, tariffUsecase, walletUsecase, log)
	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware(log))
	handler.RegisterRoutes(r)

	fmt.Println("Server started on :", cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, r)
}
