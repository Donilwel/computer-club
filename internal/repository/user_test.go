package repository

import (
	"computer-club/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	repo := NewMemoryUserRepo()

	user := &models.User{Name: "Alice", Role: models.Customer}
	err := repo.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
}

func TestGetUserByID(t *testing.T) {
	repo := NewMemoryUserRepo()

	user := &models.User{Name: "Bob", Role: models.Admin}
	repo.CreateUser(user)

	foundUser, err := repo.GetUserByID(user.ID)

	assert.NoError(t, err)
	assert.Equal(t, "Bob", foundUser.Name)
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo := NewMemoryUserRepo()

	user, err := repo.GetUserByID(99) // Не существующий ID

	assert.Error(t, err)
	assert.Nil(t, user)
}
