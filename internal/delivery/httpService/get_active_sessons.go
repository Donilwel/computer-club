package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"encoding/json"
	"net/http"
)

// GetActiveSessions возвращает активные сессии
func (h *Handler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение активных сессий")

	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при получении списка компьютеров: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	sessions := h.sessionService.GetActiveSessions(ctx)
	h.log.WithField("count", len(sessions)).Info("Активные сессии получены")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}
