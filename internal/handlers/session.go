package handlers

import (
	"computer-club/internal/middleware"
	_ "computer-club/internal/repository/models"
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

type StartSessionReq struct {
	PCNumber int   `json:"pc_number"`
	TariffID int64 `json:"tariff_id"`
}

// StartSession godoc
// @Summary      Начало сессии
// @Description  Запускает новую сессию для пользователя на указанном ПК с выбранным тарифом
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body StartSessionReq true "Данные для старта сессии"
// @Success      200 {object} models.Session
// @Failure      400 {object} string "ошибка JSON запроса"
// @Failure      401 {object} string "неправильный user_id в токене"
// @Failure      404 {object} string "Не найдено - user, computer или tariff или кошелек"
// @Failure      409 {object} string "Сессия уже активная или компьютер занят"
// @Failure      500 {object} string "внутренняя проблема поиска в базе данных или ошибка сохранения изменения транзакции базы данных или редис не создал данные"
// @Router       /session/start [post]
func (h sessionHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на начало сессии")

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.log.Error("Ошибка: user_id не найден в контексте")
		middleware.WriteError(w, http.StatusUnauthorized, errors.ErrWrongIDFromJWT.Error())
		return
	}

	var req StartSessionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()
	session, err := h.sessionService.StartSession(ctx, userID, req.PCNumber, req.TariffID)
	if err != nil {
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(session); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}

type EndSessionRequest struct {
	SessionID int64 `json:"session_id" example:"123"`
}

// EndSession godoc
// @Summary      Завершение сессии
// @Description  Завершает активную сессию по её идентификатору
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body EndSessionRequest true "Данные для завершения сессии"
// @Success      200 {object} string "message": "Session ended successfully"
// @Failure      400 {object} string "Некорректный запрос (например, session_id == 0)"
// @Failure		 409 {object} string "Сессия уже завершена или компьютер уже свободный"
// @Failure      404 {object} string "Сессия или компьютер не найдены"
// @Failure      500 {object} string "Внутренняя ошибка сервера при закрытии сессии или транзакции базы данных или редис не удалил данные"
// @Router       /session/end [post]
func (h sessionHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на завершение сессии")

	var req EndSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()

	h.log.WithField("session_id", req.SessionID).Info("Попытка завершить сессию")

	if req.SessionID == 0 {
		h.log.Error("session_id == 0, отклоняем запрос")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidSessionID.Error())
		return
	}

	err := h.sessionService.EndSession(ctx, req.SessionID)
	if err != nil {
		switch err {
		case errors.ErrUserNotFound, errors.ErrComputerNotFound, errors.ErrTariffNotFound:
			middleware.WriteError(w, http.StatusNotFound, err.Error())
		case errors.ErrStatusSessionAlreadyFinished, errors.ErrComputerAlreadyFree:
			middleware.WriteError(w, http.StatusConflict, err.Error())
		default:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		h.log.WithError(err).Error("Ошибка завершения сессии")
		return
	}

	h.log.WithField("session_id", req.SessionID).Info("Сессия успешно завершена")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Session ended successfully"}); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}

type ActiveSessionResponse struct {
	SessionID int64  `json:"session_id" example:"123"`
	UserID    int64  `json:"user_id" example:"456"`
	PCNumber  int    `json:"pc_number" example:"7"`
	TariffID  int64  `json:"tariff_id" example:"3"`
	Status    string `json:"status" example:"active"`
}

// GetActiveSessions godoc
// @Summary      Получение активных сессий
// @Description  Возвращает список всех активных сессий (только для администраторов)
// @Tags         sessions
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} ActiveSessionResponse
// @Failure      403 {object} string "ошибка доступа к текущему запросу. необходима роль пользователя: админ"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /sessions/active [get]
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(sessions); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}
