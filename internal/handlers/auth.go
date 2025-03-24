package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/models"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService usecase.AuthService
	log         *logrus.Logger
}

func NewAuthHandler(authService usecase.AuthService, log *logrus.Logger) AuthHandler {
	return &authHandler{authService: authService, log: log}
}

type RegisterUserRequest struct {
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"secret123"`
	Role     string `json:"role" example:"customer"`
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя с указанными данными
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterUserRequest true "Данные для регистрации"
// @Success      200 {object} models.User "Пользователь успешно зарегистрирован"
// @Failure      400 {object} string "Некорректные данные (пустое имя, email, короткий пароль, неверная роль)"
// @Failure      409 {object} string "Пользователь уже существует"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на регистрацию пользователя")

	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.WithError(err).Error("Ошибка декодирования JSON")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()

	role := models.UserRole(req.Role)
	if role != models.Admin && role != models.Customer {
		h.log.WithField("role", req.Role).Error("Неверная роль")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidRole.Error())
		return
	}

	user, err := h.authService.Register(ctx, req.Email, req.Password, role)
	if err != nil {
		switch err {
		case errors.ErrHashedPassword, errors.ErrRegistration:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		case errors.ErrUserAlreadyExists, errors.ErrUsernameTaken:
			middleware.WriteError(w, http.StatusConflict, err.Error())
		case errors.ErrNameEmpty, errors.ErrEmailEmpty, errors.ErrPasswordEmpty, errors.ErrPasswordTooShort:
			middleware.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			middleware.WriteError(w, http.StatusInternalServerError, errors.ErrUnexpected.Error())
		}
		h.log.WithError(err).Error("Ошибка при регистрации пользователя")
		return
	}

	h.log.WithFields(logrus.Fields{
		"user_id": user.ID,
		"name":    user.Name,
		"role":    user.Role,
	}).Info("Пользователь зарегистрирован")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
		return
	}
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Проверяет email и пароль, возвращает JWT-токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Данные для входа"
// @Success      200 {object} LoginResponse "Успешный вход"
// @Failure      400 {object} string "Некорректный запрос (невалидный JSON)"
// @Failure      401 {object} string "Ошибка авторизации - неверный email или пароль"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /auth/login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()

	token, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
		return
	}
}
