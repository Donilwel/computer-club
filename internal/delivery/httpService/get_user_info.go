package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"encoding/json"
	"net/http"
)

// GetUserInfo возвращает информацию о пользователе
func (h *Handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на получение информации о пользователе")

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.log.Error("Ошибка: user_id не найден в контексте")
		middleware.WriteError(w, http.StatusUnauthorized, errors.ErrWrongIDFromJWT.Error())
		return
	}
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		h.log.Error("Ошибка в поиске пользователя в базе данных")
		middleware.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	//тут нужно в будущем минуты потраченные количество сессий и просто баланс
	h.log.Info(w, http.StatusOK, "Получена информация о пользователе")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
