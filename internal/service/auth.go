package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"

	"computer-club/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo        domain.AuthRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(repo domain.AuthRepository, jwtSecret string, expiry time.Duration) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, tokenExpiry: expiry}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Ошибка хеширования пароля:", err)
		return "", err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Println("Ошибка создания пользователя:", err)
		return "", err
	}

	return userID, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("пользователь не найден")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("неверный пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
