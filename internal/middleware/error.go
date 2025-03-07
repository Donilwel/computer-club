package middleware

import (
	"computer-club/internal/errors"
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse - структура для возврата JSON-ошибки
type errorResponse struct {
	Message string `json:"error"`
}

// WriteError - функция для отправки JSON-ошибок
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Message: message})
}

// ErrorHandler - middleware для обработки паники
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				WriteError(w, http.StatusInternalServerError, errors.ErrUnexpected.Error())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
