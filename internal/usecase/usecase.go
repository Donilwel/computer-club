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
	GetComputersStatus() ([]models.Computer, error)
}

type ClubUsecase struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	computerRepo repository.ComputerRepository
}

func NewClubUsecase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, computerRepo repository.ComputerRepository) *ClubUsecase {
	return &ClubUsecase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		computerRepo: computerRepo,
	}
}

// RegisterUser создает нового пользователя
func (u *ClubUsecase) RegisterUser(name string, role models.UserRole) (*models.User, error) {
	user := &models.User{Name: name, Role: role}
	err := u.userRepo.CreateUser(user)
	return user, err
}

// StartSession начинает новую игровую сессию
func (u *ClubUsecase) StartSession(userID int64, pcNumber int) (*models.Session, error) {
	// Проверяем, существует ли пользователь
	_, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return u.sessionRepo.StartSession(userID, pcNumber)
}

// EndSession завершает игровую сессию
func (u *ClubUsecase) EndSession(sessionID int64) error {
	return u.sessionRepo.EndSession(sessionID)
}

// GetActiveSessions возвращает список активных сессий
func (u *ClubUsecase) GetActiveSessions() []*models.Session {
	return u.sessionRepo.GetActiveSessions()
}

func (u *ClubUsecase) GetComputersStatus() ([]models.Computer, error) {
	return u.computerRepo.GetComputers()
}
