package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"fmt"
)

type ClubService interface {
	RegisterUser(name string, role models.UserRole) (*models.User, error)
	StartSession(userID int64, pcNumber int) (*models.Session, error)
	EndSession(sessionID int64) error
	GetActiveSessions() []*models.Session
}

type clubUsecase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewClubUsecase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) ClubService {
	return &clubUsecase{userRepo: userRepo, sessionRepo: sessionRepo}
}

// RegisterUser создает нового пользователя
func (u *clubUsecase) RegisterUser(name string, role models.UserRole) (*models.User, error) {
	user := &models.User{Name: name, Role: role}
	err := u.userRepo.CreateUser(user)
	return user, err
}

// StartSession начинает новую игровую сессию
func (u *clubUsecase) StartSession(userID int64, pcNumber int) (*models.Session, error) {
	// Проверяем, существует ли пользователь
	_, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return u.sessionRepo.StartSession(userID, pcNumber)
}

// EndSession завершает игровую сессию
func (u *clubUsecase) EndSession(sessionID int64) error {
	return u.sessionRepo.EndSession(sessionID)
}

// GetActiveSessions возвращает список активных сессий
func (u *clubUsecase) GetActiveSessions() []*models.Session {
	return u.sessionRepo.GetActiveSessions()
}
