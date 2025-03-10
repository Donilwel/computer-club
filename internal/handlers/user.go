package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/repository/models"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	InfoUser(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService usecase.UserService
	log         *logrus.Logger
}

func NewUserHandler(userService usecase.UserService, log *logrus.Logger) UserHandler {
	return &userHandler{userService: userService, log: log}
}

func (h userHandler) InfoUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение информации о пользователе")

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.log.Error("Ошибка: user_id не найден в контексте")
		middleware.WriteError(w, http.StatusUnauthorized, errors.ErrWrongIDFromJWT.Error())
		return
	}
	user, err := h.userService.GetUserByID(ctx, userID)
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

func (h userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на регистрацию пользователя")

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
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

	user, err := h.userService.RegisterUser(ctx, req.Name, req.Email, req.Password, role)
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
	json.NewEncoder(w).Encode(user)
}

func (h userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
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
