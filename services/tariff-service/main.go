package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "services/config"

    "github.com/go-chi/chi/v5"
)

func main() {
    cfg := config.LoadConfig()

    r := chi.NewRouter()

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    srv := &http.Server{
        Addr:    ":" + cfg.Server.Port,
        Handler: r,
    }

    fmt.Println("Service started on port", cfg.Server.Port)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatal("Server error:", err)
    }
}