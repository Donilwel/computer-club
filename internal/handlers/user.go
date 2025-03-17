package handlers

import (
	"computer-club/internal/middleware"
	"computer-club/internal/repository/models"
	"computer-club/internal/usecase"
	"computer-club/pkg/errors"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	InfoUser(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	userService   usecase.UserService
	walletService usecase.WalletService
	log           *logrus.Logger
}

func NewUserHandler(userService usecase.UserService, walletService usecase.WalletService, log *logrus.Logger) UserHandler {
	return &userHandler{userService: userService, walletService: walletService, log: log}
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
// @Router       /info [get]
func (h userHandler) InfoUser(w http.ResponseWriter, r *http.Request) {
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
	}
}

type RegisterUserRequest struct {
	Name     string `json:"name" example:"Иван Иванов"`
	Email    string `json:"email" example:"ivan@example.com"`
	Password string `json:"password" example:"secret123"`
	Role     string `json:"role" example:"customer"`
}

// RegisterUser godoc
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя с указанными данными
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body RegisterUserRequest true "Данные для регистрации"
// @Success      200 {object} models.User "Пользователь успешно зарегистрирован"
// @Failure      400 {object} string "Некорректные данные (пустое имя, email, короткий пароль, неверная роль)"
// @Failure      409 {object} string "Пользователь уже существует"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /register [post]
func (h userHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
	}
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"secret"`
}

type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// LoginUser godoc
// @Summary      Авторизация пользователя
// @Description  Проверяет email и пароль, возвращает JWT-токен
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Данные для входа"
// @Success      200 {object} LoginResponse "Успешный вход"
// @Failure      400 {object} string "Некорректный запрос (невалидный JSON)"
// @Failure      401 {object} string "Ошибка авторизации - неверный email или пароль"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /login [post]
func (h userHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, errors.ErrJSONRequest.Error())
		return
	}
	defer r.Body.Close()

	token, err := h.userService.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		h.log.WithError(err).Error("Ошибка при кодировании ответа JSON")
		middleware.WriteError(w, http.StatusInternalServerError, errors.ErrCodingaData.Error())
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
func (h userHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.log.Info("Запрос на получение списка пользователей")

	role, ok := r.Context().Value("role").(string)
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
func (h userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
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
	}
}
