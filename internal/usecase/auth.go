package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"computer-club/pkg/JWT"
	"computer-club/pkg/errors"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email, password string, role models.UserRole) (*models.User, error)
}

type AuthUsecase struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
}

func (u *AuthUsecase) Register(ctx context.Context, email, password string, role models.UserRole) (*models.User, error) {
	// Проверки на пустые поля
	if email == "" {
		return nil, errors.ErrEmailEmpty
	}
	if password == "" {
		return nil, errors.ErrPasswordEmpty
	}
	if len(password) < 6 {
		return nil, errors.ErrPasswordTooShort
	}

	// Если роль не передана — ставим Customer по умолчанию
	if role == "" {
		role = models.Customer
	}

	// Проверяем, существует ли пользователь с таким email
	existingUser, _ := u.userRepo.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrHashedPassword
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     string(role),
	}
	tx := u.userRepo.BeginTransaction(ctx)
	defer tx.Rollback()

	if err := u.userRepo.CreateUser(ctx, tx, user); err != nil {
		return nil, err
	}

	wallet := &models.Wallet{
		UserID:  user.ID,
		Balance: 0.0,
	}

	if err := u.walletRepo.CreateWallet(ctx, tx, wallet); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	token, err := JWT.GenerateJWT(user)
	if err != nil {
		return "", errors.ErrTokenGeneration
	}

	return token, nil
}

func NewAuthUsecase(userRepo repository.UserRepository, walletRepo repository.WalletRepository) AuthService {
	return &AuthUsecase{userRepo: userRepo, walletRepo: walletRepo}
}
