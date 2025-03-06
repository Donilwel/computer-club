package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartSession(t *testing.T) {
	repo := NewMemorySessionRepo()

	session, err := repo.StartSession(1, 5)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), session.ID)
	assert.Equal(t, 5, session.PCNumber)
}

func TestEndSession(t *testing.T) {
	repo := NewMemorySessionRepo()

	session, _ := repo.StartSession(1, 3)
	err := repo.EndSession(session.ID)

	assert.NoError(t, err)
}

func TestEndSession_InvalidSession(t *testing.T) {
	repo := NewMemorySessionRepo()

	err := repo.EndSession(99) // Не существующий ID

	assert.Error(t, err)
}
