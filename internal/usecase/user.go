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

type UserService interface {
	RegisterUser(name, email, password string, role models.UserRole) (*models.User, error)
	LoginUser(name string, password string) (string, error)
	GetUserByEmail(email string) (*models.User, error)
}

type UserUsercase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) *UserUsercase {
	return &UserUsercase{userRepo: userRepo}
}

func (u *UserUsercase) RegisterUser(name, email, password string, role models.UserRole) (*models.User, error) {
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

func (u *UserUsercase) LoginUser(email string, password string) (string, error) {
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

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func generateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен живёт 24 часа
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func (u *UserUsercase) GetUserByEmail(email string) (*models.User, error) {
	return u.userRepo.GetUserByEmail(email)
}
