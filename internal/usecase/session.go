package usecase

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"computer-club/internal/repository"
)

type SessionService interface {
	StartSession(userID int64, pcNumber int) (*models.Session, error)
	EndSession(sessionID int64) error
	GetActiveSessions() []*models.Session
}

type SessionUsecase struct {
	sessionRepository repository.SessionRepository
	userRepo          repository.UserRepository
}

func NewSessionUsecase(sessionRepository repository.SessionRepository, userRepo repository.UserRepository) *SessionUsecase {
	return &SessionUsecase{sessionRepository: sessionRepository, userRepo: userRepo}
}

func (u *SessionUsecase) StartSession(userID int64, pcNumber int) (*models.Session, error) {
	_, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	return u.sessionRepository.StartSession(userID, pcNumber)
}

func (u *SessionUsecase) EndSession(sessionID int64) error {
	return u.sessionRepository.EndSession(sessionID)
}

func (u *SessionUsecase) GetActiveSessions() []*models.Session {
	return u.sessionRepository.GetActiveSessions()
}
