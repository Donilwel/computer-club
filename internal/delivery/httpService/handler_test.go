package httpService

import (
	"bytes"
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"computer-club/internal/usecase"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func setupTestServer() *httptest.Server {
	userRepo := repository.NewMemoryUserRepo()
	sessionRepo := repository.NewMemorySessionRepo()
	service := usecase.NewClubUsecase(userRepo, sessionRepo)
	handler := NewHandler(service)

	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	return httptest.NewServer(r)
}

func TestRegisterUserHandler(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	reqBody := `{"name": "Alice", "role": "customer"}`
	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewBufferString(reqBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var user models.User
	json.NewDecoder(resp.Body).Decode(&user)

	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, models.Customer, user.Role)
}

func TestStartSessionHandler(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Создаем пользователя
	http.Post(server.URL+"/register", "application/json", bytes.NewBufferString(`{"name": "Bob", "role": "customer"}`))

	reqBody := `{"user_id": 1, "pc_number": 2}`
	resp, err := http.Post(server.URL+"/session/start", "application/json", bytes.NewBufferString(reqBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestStartSessionHandler_InvalidUser(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	reqBody := `{"user_id": 99, "pc_number": 2}`
	resp, err := http.Post(server.URL+"/session/start", "application/json", bytes.NewBufferString(reqBody))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestMassiveRegisterUserAPI(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reqBody := `{"name": "` + randomString(8) + `", "role": "customer"}`
			resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewBufferString(reqBody))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}()
	}
	wg.Wait()
}
