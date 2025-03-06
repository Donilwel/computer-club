package httpService

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"computer-club/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	clubService usecase.ClubService
	log         *logrus.Logger
}

// NewHandler создает новый HTTP-обработчик
func NewHandler(clubService usecase.ClubService, log *logrus.Logger) *Handler {
	return &Handler{clubService: clubService, log: log}
}

// RegisterRoutes регистрирует эндпоинты
func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Post("/register", h.RegisterUser)
	r.Post("/session/start", h.StartSession)
	r.Post("/session/end", h.EndSession)
	r.Get("/sessions/active", h.GetActiveSessions)
	r.Get("/computers/status", h.GetComputersStatus)
}

// RegisterUser регистрирует нового пользователя
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на регистрацию пользователя")

	var req struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		writeError(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	role := models.UserRole(req.Role)
	if role != models.Admin && role != models.Customer {
		h.log.WithField("role", req.Role).Error("Неверная роль")
		writeError(w, http.StatusBadRequest, "Invalid role")
		return
	}

	user, err := h.clubService.RegisterUser(req.Name, role)
	if err != nil {
		h.log.WithError(err).Error("Ошибка при регистрации пользователя")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"name":    user.Name,
		"role":    user.Role,
	}).Info("Пользователь зарегистрирован")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// StartSession начинает новую сессию
func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на начало сессии")

	var req struct {
		UserID   int64 `json:"user_id"`
		PCNumber int   `json:"pc_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		writeError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	session, err := h.clubService.StartSession(req.UserID, req.PCNumber)
	if err != nil {
		// Проверяем тип ошибки
		switch err {
		case errors.ErrUserNotFound:
			writeError(w, http.StatusNotFound, err.Error())
		case errors.ErrSessionActive:
			writeError(w, http.StatusConflict, err.Error())
		case errors.ErrPCBusy:
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Internal server error")
		}

		h.log.WithError(err).Error("Ошибка при запуске сессии")
		return
	}

	h.log.WithFields(logrus.Fields{
		"session_id": session.ID,
		"user_id":    session.UserID,
		"pc_number":  session.PCNumber,
	}).Info("Сессия успешно запущена")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// EndSession завершает сессию
func (h *Handler) EndSession(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на завершение сессии")

	var req struct {
		SessionID int64 `json:"session_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		writeError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	h.log.WithField("session_id", req.SessionID).Info("Попытка завершить сессию")

	if req.SessionID == 0 {
		h.log.Error("session_id == 0, отклоняем запрос")
		writeError(w, http.StatusBadRequest, "Invalid session_id")
		return
	}

	err := h.clubService.EndSession(req.SessionID)
	if err != nil {
		h.log.WithError(err).Error("Ошибка завершения сессии")
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.log.WithField("session_id", req.SessionID).Info("Сессия успешно завершена")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Session ended successfully"})
}

// GetActiveSessions возвращает активные сессии
func (h *Handler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на получение активных сессий")

	sessions := h.clubService.GetActiveSessions()
	h.log.WithField("count", len(sessions)).Info("Активные сессии получены")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// GetComputersStatus возвращает статус компьютеров
func (h *Handler) GetComputersStatus(w http.ResponseWriter, r *http.Request) {
	h.log.Info("Запрос на получение статуса компьютеров")

	computers, err := h.clubService.GetComputersStatus()
	if err != nil {
		h.log.WithError(err).Error("Ошибка при получении списка компьютеров")
		writeError(w, http.StatusInternalServerError, "Ошибка при получении списка компьютеров")
		return
	}

	h.log.WithField("count", len(computers)).Info("Статус компьютеров получен")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(computers)
}
