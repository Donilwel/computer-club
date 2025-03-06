package main

import (
	"fmt"
	"net/http"

	"computer-club/config"
	"computer-club/internal/delivery/httpService"
	"computer-club/internal/repository"
	"computer-club/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.LoadConfig()

	// Подключаемся к PostgreSQL
	db := repository.NewPostgresDB(cfg)
	repository.Migrate(db)

	// Используем PostgreSQL репозитории
	userRepo := repository.NewPostgresUserRepo(db)
	sessionRepo := repository.NewPostgresSessionRepo(db)

	// Создаем бизнес-логику
	clubService := usecase.NewClubUsecase(userRepo, sessionRepo)

	// Создаем HTTP обработчик
	handler := httpService.NewHandler(clubService)

	// Настраиваем маршрутизатор
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(httpService.ErrorHandler)

	// Регистрируем API
	handler.RegisterRoutes(r)

	// Запускаем сервер
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
