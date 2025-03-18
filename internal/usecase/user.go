package usecase

import (
	"computer-club/internal/repository"
	"computer-club/internal/repository/models"
	"context"
)

type UserService interface {
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

func (u *UserUsecase) GetUsers(ctx context.Context) ([]*models.User, error) {
	return u.userRepo.GetUsers(ctx)
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
