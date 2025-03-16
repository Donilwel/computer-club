package utils

import (
    "encoding/json"
    "net/http"
)

type ErrorResponse struct {
    Message string `json:"error"`
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}
