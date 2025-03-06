package httpService

import (
	"encoding/json"
	"net/http"
	"strconv"

	"computer-club/internal/models"
	"computer-club/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	clubService usecase.ClubService
}

// NewHandler создает новый HTTP-обработчик
func NewHandler(clubService usecase.ClubService) *Handler {
	return &Handler{clubService: clubService}
}

// RegisterRoutes регистрирует эндпоинты
func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Post("/register", h.RegisterUser)
	r.Post("/session/start", h.StartSession)
	r.Post("/session/end", h.EndSession)
	r.Get("/sessions/active", h.GetActiveSessions)
}

// RegisterUser регистрирует нового пользователя
func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON request")
		return
	}

	role := models.UserRole(req.Role)
	if role != models.Admin && role != models.Customer {
		writeError(w, http.StatusBadRequest, "Invalid role")
		return
	}

	user, err := h.clubService.RegisterUser(req.Name, role)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// StartSession начинает новую сессию
func (h *Handler) StartSession(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int64 `json:"user_id"`
		PCNumber int   `json:"pc_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	session, err := h.clubService.StartSession(req.UserID, req.PCNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// EndSession завершает сессию
func (h *Handler) EndSession(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.ParseInt(r.URL.Query().Get("session_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid session_id", http.StatusBadRequest)
		return
	}

	if err := h.clubService.EndSession(sessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetActiveSessions возвращает активные сессии
func (h *Handler) GetActiveSessions(w http.ResponseWriter, r *http.Request) {
	sessions := h.clubService.GetActiveSessions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}
