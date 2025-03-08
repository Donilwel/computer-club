package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/middleware"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// StartSession начинает новую сессию
func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на начало сессии")

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.log.Error("Ошибка: user_id не найден в контексте")
		middleware.WriteError(w, http.StatusUnauthorized, errors.ErrWrongIDFromJWT.Error())
		return
	}

	var req struct {
		PCNumber int   `json:"pc_number"`
		TariffID int64 `json:"tariff_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	session, err := h.sessionService.StartSession(userID, req.PCNumber, req.TariffID)
	if err != nil {
		// Проверяем тип ошибки
		switch err {
		case errors.ErrUserNotFound, errors.ErrComputerNotFound, errors.ErrTariffNotFound:
			middleware.WriteError(w, http.StatusNotFound, err.Error())
		case errors.ErrSessionActive, errors.ErrPCBusy:
			middleware.WriteError(w, http.StatusConflict, err.Error())
		case errors.ErrCreatedSession:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		default:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		h.log.WithError(err).Error("Ошибка при запуске сессии")
		return
	}

	h.log.WithFields(logrus.Fields{
		"session_id": session.ID,
		"user_id":    session.UserID,
		"pc_number":  session.PCNumber,
		"tariff_id":  session.TariffID,
		"status":     session.Status,
	}).Info("Сессия успешно запущена")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}
