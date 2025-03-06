package repository

import (
	"computer-club/internal/models"
	"fmt"
	"sync"
	"time"
)

type SessionRepository interface {
	StartSession(userID int64, pcNumber int) (*models.Session, error)
	EndSession(sessionID int64) error
	GetActiveSessions() []*models.Session
}

type memorySessionRepo struct {
	mu       sync.Mutex
	sessions map[int64]*models.Session
	lastID   int64
}

func NewMemorySessionRepo() SessionRepository {
	return &memorySessionRepo{
		sessions: make(map[int64]*models.Session),
	}
}

func (r *memorySessionRepo) StartSession(userID int64, pcNumber int) (*models.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	session := &models.Session{
		ID:        r.lastID + 1,
		UserID:    userID,
		PCNumber:  pcNumber,
		StartTime: time.Now(),
	}
	r.sessions[session.ID] = session
	return session, nil
}

// EndSession завершает сессию
func (r *memorySessionRepo) EndSession(sessionID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	session, exists := r.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	now := time.Now()
	session.EndTime = &now
	return nil
}

// GetActiveSessions возвращает список активных сессий
func (r *memorySessionRepo) GetActiveSessions() []*models.Session {
	r.mu.Lock()
	defer r.mu.Unlock()

	var activeSessions []*models.Session
	for _, session := range r.sessions {
		if session.EndTime == nil {
			activeSessions = append(activeSessions, session)
		}
	}
	return activeSessions
}
