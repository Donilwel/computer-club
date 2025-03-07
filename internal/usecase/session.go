package usecase

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"context"
	"log"
	"time"
)

type SessionService interface {
	StartSession(userID int64, pcNumber int, tariffID int64) (*models.Session, error)
	EndSession(sessionID int64) error
	GetActiveSessions() []*models.Session
	MonitorSessions(ctx context.Context)
}

type SessionUsecase struct {
	sessionRepository repository.SessionRepository
	userRepo          repository.UserRepository
}

func NewSessionUsecase(sessionRepository repository.SessionRepository, userRepo repository.UserRepository) *SessionUsecase {
	return &SessionUsecase{sessionRepository: sessionRepository, userRepo: userRepo}
}

func (u *SessionUsecase) StartSession(userID int64, pcNumber int, tariffID int64) (*models.Session, error) {
	_, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	return u.sessionRepository.StartSession(userID, pcNumber, tariffID)
}

func (u *SessionUsecase) EndSession(sessionID int64) error {
	return u.sessionRepository.EndSession(sessionID)
}

func (u *SessionUsecase) GetActiveSessions() []*models.Session {
	return u.sessionRepository.GetActiveSessions()
}

func (u *SessionUsecase) MonitorSessions(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Остановка мониторинга сессий")
			return
		case <-ticker.C:
			u.checkAndCloseExpiredSessions()
		}
	}
}

func (u *SessionUsecase) checkAndCloseExpiredSessions() {
	sessions := u.sessionRepository.GetActiveSessions()
	now := time.Now()

	for _, session := range sessions {
		if session.EndTime.Before(now) {
			log.Printf("Завершаем сессию %d (пользователь %d)", session.ID, session.UserID)
			u.EndSession(session.ID)
		}
	}
}
