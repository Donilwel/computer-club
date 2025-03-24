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
	authHandler handlers.AuthHandler,
) {
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Get("/tariff", tariffHandler.GetTariff)
	r.Get("/tariff/{id}", tariffHandler.GetTariffByID)

	r.Group(func(protected chi.Router) {
		protected.Use(middleware.AuthMiddleware)

		protected.Post("/tariff", tariffHandler.AddTariff)
		//protected.Put("/tariff/{id}", tariffHandler.ChangeTariff)
		protected.Delete("/tariff/{id}", tariffHandler.DeleteTariff)

		protected.Get("/users", userHandler.GetUsers)
		protected.Get("/users/{id}", userHandler.GetUserByID)
		protected.Get("/users/info", userHandler.InfoUser)

		protected.Put("/pay", walletHandler.PutMoneyOnWallet)

		protected.Post("/session/start", sessionHandler.StartSession)
		protected.Get("/session/active", sessionHandler.GetActiveSessions)
		protected.Post("/session/end", sessionHandler.EndSession)

		protected.Get("/computers", computerHandler.GetComputersStatus)
		protected.Post("/computers", computerHandler.AddComputer)
		protected.Delete("/computers/{id}", computerHandler.DeleteComputer)
	})
}
