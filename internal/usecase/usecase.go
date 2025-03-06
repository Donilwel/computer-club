package usecase

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type ClubService interface {
	RegisterUser(name, email, password string, role models.UserRole) (*models.User, error)
	LoginUser(name string, password string) (string, error)
	StartSession(userID int64, pcNumber int) (*models.Session, error)
	EndSession(sessionID int64) error
	GetActiveSessions() []*models.Session
	GetComputersStatus() ([]models.Computer, error)
	GetUserByEmail(email string) (*models.User, error)
}

type ClubUsecase struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	computerRepo repository.ComputerRepository
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func (u *ClubUsecase) LoginUser(email string, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	token, err := generateJWT(user)
	if err != nil {
		return "", errors.ErrTokenGeneration
	}

	return token, nil
}

func generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен живёт 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func NewClubUsecase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, computerRepo repository.ComputerRepository) *ClubUsecase {
	return &ClubUsecase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		computerRepo: computerRepo,
	}
}

// RegisterUser создает нового пользователя
func (u *ClubUsecase) RegisterUser(name, email, password string, role models.UserRole) (*models.User, error) {
	existingUser, _ := u.userRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	existingUserByName, _ := u.userRepo.GetUserByEmail(email)
	if existingUserByName != nil {
		return nil, errors.ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     string(role),
	}

	if err := u.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// StartSession начинает новую игровую сессию
func (u *ClubUsecase) StartSession(userID int64, pcNumber int) (*models.Session, error) {
	// Проверяем, существует ли пользователь
	_, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
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

func (u *ClubUsecase) GetUserByEmail(email string) (*models.User, error) {
	return u.userRepo.GetUserByEmail(email)
}
