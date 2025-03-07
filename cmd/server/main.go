package main

import (
	"computer-club/internal/logger"
	"computer-club/internal/middleware"
	"computer-club/internal/usecase"
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
	userRepo := repository.NewPostgresUserRepo(db)
	sessionRepo := repository.NewPostgresSessionRepo(db, redisClient)
	computerRepo := repository.NewComputerRepository(db)
	tariffRepo := repository.NewTariffRepositoryPostgres(db)

	userUsecase := usecase.NewUserUsecase(userRepo)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo)
	computerUsecase := usecase.NewComputerUsecase(computerRepo)
	tariffUsecase := usecase.NewTariffUsecase(tariffRepo)
	// Запускаем HTTP сервер
	handler := httpService.NewHandler(userUsecase, computerUsecase, sessionUsecase, tariffUsecase, log)
	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware(log))
	handler.RegisterRoutes(r)

	fmt.Println("Server started on :", cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, r)
}
