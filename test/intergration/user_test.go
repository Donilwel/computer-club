package httpService_test

import (
	"bytes"
	"computer-club/internal/handlers"
	"computer-club/internal/models"
	"computer-club/pkg/errors"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Подготовка мок-сервиса для регистрации
type mockAuthService struct{}

func (m *mockAuthService) Register(ctx context.Context, email, password string, role models.UserRole) (*models.User, error) {
	if email == "test@example.com" {
		return nil, errors.ErrUserAlreadyExists
	}
	return &models.User{ID: 1, Email: email, Role: string(role)}, nil
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, error) {
	return "token", nil
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "User already exists",
			requestBody: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
				"role":     "customer",
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   map[string]interface{}{"error": "Пользователь уже существует"},
		},
		{
			name: "Successful registration",
			requestBody: map[string]string{
				"email":    "newuser@example.com",
				"password": "password123",
				"role":     "customer",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"id": float64(1), "email": "newuser@example.com", "role": "customer"},
		},
		{
			name: "Missing email",
			requestBody: map[string]string{
				"email":    "",
				"password": "password123",
				"role":     "customer",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "Некорректные данные (пустое имя, email, короткий пароль, неверная роль)"},
		},
		{
			name: "Invalid role",
			requestBody: map[string]string{
				"email":    "invalidrole@example.com",
				"password": "password123",
				"role":     "invalidrole",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "некорректная роль"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			log := logrus.New()
			authHandler := handlers.NewAuthHandler(&mockAuthService{}, log)
			r.Post("/auth/register", authHandler.Register)

			// Подготовка запроса
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Ответ
			w := httptest.NewRecorder()

			// Выполнение запроса
			r.ServeHTTP(w, req)

			// Проверка статуса ответа
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Проверка тела ответа
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.Nil(t, err)

			// Сравнение с ожидаемым ответом
			assert.Equal(t, tt.expectedBody, response)
		})
	}
}
