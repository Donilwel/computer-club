package httpService

import (
	"encoding/json"
	"log"
	"net/http"
)

// errorResponse - структура для возврата JSON-ошибки
type errorResponse struct {
	Message string `json:"error"`
}

// writeError - функция для отправки JSON-ошибок
func writeError(w http.ResponseWriter, statusCode int, message string) {
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
				writeError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
