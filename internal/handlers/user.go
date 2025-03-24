package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/models"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserHandler interface {
	InfoUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService usecase.UserService
	log         *logrus.Logger
}

func NewUserHandler(userService usecase.UserService, log *logrus.Logger) UserHandler {
	return &userHandler{userService: userService, log: log}
}

type UserInfoResponse struct {
	User         *models.User          `json:"user"`
	Balance      float64               `json:"balance" example:"100.50"`
	Transactions []*models.Transaction `json:"transactions"`
}

// InfoUser godoc
// @Summary      Получение информации о пользователе
// @Description  Возвращает информацию о пользователе, баланс и транзакции
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} UserInfoResponse
// @Failure      401 {object} string "Ошибка авторизации - user_id не найден"
// @Failure      400 {object} string "Некорректный user_id"
// @Failure      401 {object} string "Кошелек не найден или пользователь"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /users/info [get]
func (h *userHandler) InfoUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение информации о пользователе")

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.log.Error("Ошибка: user_id не найден в контексте")
		middleware.WriteError(w, http.StatusUnauthorized, errors.ErrWrongIDFromJWT.Error())
		return
	}
	user, balance, transactions, err := h.userService.GetInfoUser(ctx, userID)
	if err != nil {
		switch err {
		case errors.ErrFindUser, errors.ErrCheckBalance:
			middleware.WriteError(w, http.StatusNotFound, err.Error())
		case errors.ErrCheckTransaction:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		default:
			middleware.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		h.log.Error(err)
		return
	}

	response := UserInfoResponse{
		User:         user,
		Balance:      balance,
		Transactions: transactions,
	}

	h.log.Info(w, http.StatusOK, "Получена информация о пользователе")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
		return
	}
}

// GetUsers godoc
// @Summary      Получение списка пользователей
// @Description      Получение списка пользователей (может только админ)
// @Tags 		users
// @Produce      json
// @Security     BearerAuth
// @Success 	 200 {object} []models.User "Список пользователей успешно показан"
// @Failure		 403 {object} string "Недостаточно прав"
// @Failure		 404 {object} string "Список пользователей пуст"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router 		/users [get]
func (h *userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение списка пользователей")

	role, ok := ctx.Value("role").(string)
	if !ok || role != "admin" {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при получении списка пользователей: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	users, err := h.userService.GetUsers(ctx)
	if err != nil {
		h.log.WithError(err).Error("Ошибка при получении списка пользователей")
		if err == errors.ErrZeroUsers {
			middleware.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrFindUsers.Error())
		return
	}

	h.log.WithField("count", len(users)).Info("Список пользователей получены")

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
		return
	}
}

// GetUserByID
// @Summary  		Получение пользователя по ID
// @Description      Получение пользователя по ID (может только админ)
// @Tags 		 users
// @Produce      json
// @Security     BearerAuth
// @Success 	 200 {object} models.User "Пользователь успешно показан"
// @Failure 	 400 {object} string "Некорректный id"
// @Failure		 403 {object} string "Недостаточно прав"
// @Failure		 404 {object} string "Пользователь не найден"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router 		/users/{id} [get]
func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	role, ok := r.Context().Value("role").(string)
	if !ok || role != "admin" {
		h.log.WithError(errors.ErrForbidden).Error("Ошибка при получении списка пользователей: недостаточно прав")
		middleware.WriteError(w, http.StatusForbidden, errors.ErrForbidden.Error())
		return
	}

	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.WithError(err).Error("Некорректный ID пользователя")
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrInvalidUserID.Error())
		return
	}

	h.log.Info("Запрос на получение пользователя по id")
	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		h.log.WithError(err).Error("Ошибка при поиске пользователя по ID")
		middleware.WriteError(w, http.StatusNotFound, errors.ErrFindUsers.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
		return
	}
}
