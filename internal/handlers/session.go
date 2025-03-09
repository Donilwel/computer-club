package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SessionHandler interface {
	StartSession(w http.ResponseWriter, r *http.Request)
	EndSession(w http.ResponseWriter, r *http.Request)
	GetActiveSessions(w http.ResponseWriter, r *http.Request)
}

type sessionHandler struct {
	sessionService usecase.SessionService
	log            *logrus.Logger
}

func NewSessionHandler(sessionService usecase.SessionService, log *logrus.Logger) SessionHandler {
	return sessionHandler{sessionService: sessionService, log: log}
}

// StartSession начинает новую сессию
func (h sessionHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	session, err := h.sessionService.StartSession(ctx, userID, req.PCNumber, req.TariffID)
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

// EndSession завершает активную сессию
func (h sessionHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	err := h.sessionService.EndSession(ctx, req.SessionID)
	if err != nil {
		h.log.WithError(err).Error("Ошибка завершения сессии")
		middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.WithField("session_id", req.SessionID).Info("Сессия успешно завершена")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Session ended successfully"})
}

// GetActiveSessions выводит список активных сессий
func (h sessionHandler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
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
