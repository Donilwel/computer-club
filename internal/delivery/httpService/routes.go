package httpService

import (
	"computer-club/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes регистрирует эндпоинты
func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Post("/register", h.RegisterUser)
	r.Post("/login", h.LoginUser)
	r.Get("/tariff", h.GetTariff)
	r.Get("/tariff/{id}", h.GetTariffByID)
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)
		protected.Get("/info", h.GetUserInfo)
		protected.Post("/session/start", h.StartSession)
		protected.Post("/session/end", h.EndSession)
		protected.Get("/sessions/active", h.GetActiveSessions)
		protected.Get("/computers/status", h.GetComputersStatus)
	})
}
