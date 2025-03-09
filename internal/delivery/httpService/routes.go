package httpService

import (
	"computer-club/internal/handlers"
	"computer-club/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes регистрирует эндпоинты
func RegisterRoutes(
	r *chi.Mux,
	userHandler handlers.UserHandler,
	tariffHandler handlers.TariffHandler,
	sessionHandler handlers.SessionHandler,
	walletHandler handlers.WalletHandler,
	computerHandler handlers.ComputerHandler,
) {
	r.Post("/register", userHandler.RegisterUser)
	r.Post("/login", userHandler.LoginUser)

	r.Get("/tariff", tariffHandler.GetTariff)
	r.Get("/tariff/{id}", tariffHandler.GetTariffByID)

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)

		protected.Get("/info", userHandler.InfoUser)
		protected.Post("/session/start", sessionHandler.StartSession)
		protected.Post("/session/end", sessionHandler.EndSession)
		protected.Put("/pay", walletHandler.PutMoneyOnWallet)
		protected.Get("/sessions/active", sessionHandler.GetActiveSessions)
		protected.Get("/computers/status", computerHandler.GetComputersStatus)
	})
}
