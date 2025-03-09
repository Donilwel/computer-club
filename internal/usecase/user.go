package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"computer-club/pkg/errors"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type UserService interface {
	RegisterUser(ctx context.Context, name, email, password string, role models.UserRole) (*models.User, error)
	LoginUser(ctx context.Context, name string, password string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
}

type UserUsecase struct {
	userRepo      repository.UserRepository
	walletService WalletService
}

func NewUserUsecase(userRepo repository.UserRepository,
	walletService WalletService) UserService {
	return &UserUsecase{userRepo: userRepo, walletService: walletService}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, name, email, password string, role models.UserRole) (*models.User, error) {
	// Проверки на пустые поля
	if name == "" {
		return nil, errors.ErrNameEmpty
	}
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

	// Проверяем, существует ли пользователь с таким name
	existingUserByName, _ := u.userRepo.GetUserByName(ctx, name)
	if existingUserByName != nil {
		return nil, errors.ErrUsernameTaken
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrHashedPassword
	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     string(role),
	}

	// Сохраняем пользователя в БД
	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	if err := u.walletService.CreateWallet(ctx, user.ID); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) LoginUser(ctx context.Context, email string, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
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

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}
