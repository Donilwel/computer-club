package main

import (
	"github.com/go-chi/chi/v5"
	"services/handlers"
)

func RegisterRoutes(r *chi.Mux) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/api/items", func(api chi.Router) {
		api.Get("/", handlers.GetItems)
		api.Get("/{id}", handlers.GetItemByID)
		api.Post("/", handlers.CreateItem)
		api.Put("/{id}", handlers.UpdateItem)
		api.Delete("/{id}", handlers.DeleteItem)
	})
}
