package auth

import (
	"context"
	"errors"
	"log"
)

// AuthService интерфейс
type AuthService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (bool, error)
}

// authServiceImpl — реализация интерфейса
type authServiceImpl struct {
	repo AuthRepository
}

// NewAuthService создает новый сервис
func NewAuthService(repo AuthRepository) AuthService {
	return &authServiceImpl{repo: repo}
}

func (s *authServiceImpl) Register(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password required")
	}
	userID, err := s.repo.CreateUser(ctx, username, password)
	if err != nil {
		return "", err
	}
	log.Printf("User registered: %s", username)
	return userID, nil
}

func (s *authServiceImpl) Login(ctx context.Context, username, password string) (string, error) {
	token, err := s.repo.CheckUser(ctx, username, password)
	if err != nil {
		return "", err
	}
	log.Printf("User logged in: %s", username)
	return token, nil
}

func (s *authServiceImpl) ValidateToken(ctx context.Context, token string) (bool, error) {
	return s.repo.ValidateToken(ctx, token)
}
