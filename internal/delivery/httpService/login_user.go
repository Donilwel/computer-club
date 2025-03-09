package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"encoding/json"
	"net/http"
)

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()

	// Вызываем usecase для логина
	token, err := h.userService.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Отправляем токен
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
