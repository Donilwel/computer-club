package handlers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"computer-club/internal/di"
	"computer-club/internal/repository/models"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var container *di.Container

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %v", err)
	}

	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=test",
		"POSTGRES_PASSWORD=test",
		"POSTGRES_DB=testdb",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %v", err)
	}

	// Ждем, пока БД станет доступной
	time.Sleep(time.Second * 5)

	container = di.NewContainer()
	if container == nil {
		log.Fatal("Failed to initialize DI container")
	}
	if container.DB == nil {
		log.Fatal("Database connection is nil")
	}
	if container.UserHandler == nil {
		log.Fatal("UserHandler is nil")
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}
	os.Exit(code)
}

func resetDatabase() {
	if container.DB == nil {
		log.Fatal("DB is nil in resetDatabase")
	}
	container.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

func TestRegisterUser_Success(t *testing.T) {
	resetDatabase()

	body, _ := json.Marshal(map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password",
		"role":     "customer",
	})

	r := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	w := httptest.NewRecorder()
	container.UserHandler.RegisterUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	var respUser models.User
	json.NewDecoder(res.Body).Decode(&respUser)
	assert.Equal(t, "Test User", respUser.Name)
	assert.Equal(t, "test@example.com", respUser.Email)
}

func TestLoginUser_Success(t *testing.T) {
	resetDatabase()
	TestRegisterUser_Success(t)

	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password",
	})

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	container.UserHandler.LoginUser(w, r)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	var resp map[string]string
	json.NewDecoder(res.Body).Decode(&resp)
	assert.NotEmpty(t, resp["token"])
}
