package usecase

import (
	"computer-club/internal/repository"
	"computer-club/internal/repository/models"
	"computer-club/pkg/JWT"
	"computer-club/pkg/errors"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, name, email, password string, role models.UserRole) (*models.User, error)
	LoginUser(ctx context.Context, name string, password string) (string, error)
	GetInfoUser(ctx context.Context, userID int64) (*models.User, float64, []*models.Transaction, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
}

type UserUsecase struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewUserUsecase(userRepo repository.UserRepository,
	walletRepo repository.WalletRepository) UserService {
	return &UserUsecase{userRepo: userRepo, walletRepo: walletRepo}
}

func (u *UserUsecase) RegisterUser(ctx context.Context,
	name, email, password string, role models.UserRole) (*models.User, error) {
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

func (u *UserUsecase) GetUsers(ctx context.Context) ([]*models.User, error) {
	return u.userRepo.GetUsers(ctx)
}

func (u *UserUsecase) LoginUser(ctx context.Context, email string, password string) (string, error) {
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

func (u *UserUsecase) GetInfoUser(ctx context.Context, userID int64) (*models.User, float64, []*models.Transaction, error) {
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, 0.0, nil, err
	}

	balance, err := u.walletRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, 0.0, nil, err
	}

	transactions, err := u.walletRepo.GetTransactions(ctx, userID)
	if err != nil {
		return nil, 0.0, nil, err
	}
	return user, balance, transactions, nil
}

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.userRepo.GetUserByEmail(ctx, email)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	return u.userRepo.GetUserByID(ctx, id)
}
