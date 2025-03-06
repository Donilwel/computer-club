package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	userRepo := repository.NewMemoryUserRepo()
	sessionRepo := repository.NewMemorySessionRepo()
	service := NewClubUsecase(userRepo, sessionRepo)

	user, err := service.RegisterUser("Alice", models.Customer)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
	assert.Equal(t, models.Customer, user.Role)
}

func TestStartSession(t *testing.T) {
	userRepo := repository.NewMemoryUserRepo()
	sessionRepo := repository.NewMemorySessionRepo()
	service := NewClubUsecase(userRepo, sessionRepo)

	user, _ := service.RegisterUser("Bob", models.Customer)
	session, err := service.StartSession(user.ID, 1)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, user.ID, session.UserID)
	assert.Equal(t, 1, session.PCNumber)
}

func TestStartSession_InvalidUser(t *testing.T) {
	userRepo := repository.NewMemoryUserRepo()
	sessionRepo := repository.NewMemorySessionRepo()
	service := NewClubUsecase(userRepo, sessionRepo)

	session, err := service.StartSession(99, 1) // Не существующий user ID
	assert.Error(t, err)
	assert.Nil(t, session)
}

func TestEndSession(t *testing.T) {
	userRepo := repository.NewMemoryUserRepo()
	sessionRepo := repository.NewMemorySessionRepo()
	service := NewClubUsecase(userRepo, sessionRepo)

	user, _ := service.RegisterUser("Charlie", models.Customer)
	session, _ := service.StartSession(user.ID, 2)

	err := service.EndSession(session.ID)
	assert.NoError(t, err)
}

func TestEndSession_InvalidSession(t *testing.T) {
	userRepo := repository.NewMemoryUserRepo()
	sessionRepo := repository.NewMemorySessionRepo()
	service := NewClubUsecase(userRepo, sessionRepo)

	err := service.EndSession(999) // Не существующий session ID
	assert.Error(t, err)
}
