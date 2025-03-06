package main

import (
	"fmt"
	"net/http"

	"computer-club/internal/delivery/httpService"
	"computer-club/internal/repository"
	"computer-club/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Инициализируем репозитории
	userRepo := repository.NewMemoryUserRepo()
	sessionRepo := repository.NewMemorySessionRepo()

	// Создаем бизнес-логику
	clubService := usecase.NewClubUsecase(userRepo, sessionRepo)

	// Создаем HTTP обработчик
	handler := httpService.NewHandler(clubService)

	// Настраиваем маршрутизатор
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(httpService.ErrorHandler) // Добавляем обработку ошибок

	// Регистрируем эндпоинты
	handler.RegisterRoutes(r)

	// Запускаем сервер
	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
