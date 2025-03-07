package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"encoding/json"
	"net/http"
)

// EndSession завершает сессию
func (h *Handler) EndSession(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на завершение сессии")

	var req struct {
		SessionID int64 `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}

	h.log.WithField("session_id", req.SessionID).Info("Попытка завершить сессию")

	if req.SessionID == 0 {
		h.log.Error("session_id == 0, отклоняем запрос")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidSessionID.Error())
		return
	}

	err := h.sessionService.EndSession(req.SessionID)
	if err != nil {
		h.log.WithError(err).Error("Ошибка завершения сессии")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.WithField("session_id", req.SessionID).Info("Сессия успешно завершена")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Session ended successfully"})
}
