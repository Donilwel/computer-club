package usecase

import (
	"computer-club/internal/repository"
	models2 "computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
)

type WalletService interface {
	PutMoneyOnWallet(ctx context.Context, userID int64, amount float64) (*models2.Transaction, error)
	Withdraw(ctx context.Context, userID int64, amount float64) error
	GetBalance(ctx context.Context, userID int64) (float64, error)
	GetTransactions(ctx context.Context, userID int64) ([]*models2.Transaction, error)
}

type WalletUsecase struct {
	walletRepo repository.WalletRepository
	tariffRepo repository.TariffRepository
	userRepo   repository.UserRepository
}

func NewWalletUsecase(walletRepo repository.WalletRepository,
	tariffRepo repository.TariffRepository,
	userRepo repository.UserRepository) WalletService {
	return &WalletUsecase{walletRepo: walletRepo,
		tariffRepo: tariffRepo,
		userRepo:   userRepo}
}

func (u *WalletUsecase) PutMoneyOnWallet(ctx context.Context, userID int64, amount float64) (*models2.Transaction, error) {
	if amount <= 0 {
		return nil, errors.ErrInvalidAmount
	}

	if _, err := u.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	if _, err := u.walletRepo.GetBalance(ctx, userID); err != nil {
		return nil, err
	}

	tx := u.walletRepo.BeginTransaction(ctx)
	defer tx.Rollback()

	if err := u.walletRepo.Deposit(ctx, tx, userID, amount); err != nil {
		return nil, err
	}

	transaction, err := u.walletRepo.CreateTransaction(ctx, tx, userID, amount, string(models2.Add), nil)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (u *WalletUsecase) Withdraw(ctx context.Context, userID int64, amount float64) error {
	if amount <= 0 {
		return errors.ErrInvalidAmount
	}
	balance, err := u.walletRepo.GetBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.ErrInsufficientFunds
	}
	return u.walletRepo.Withdraw(ctx, nil, userID, amount)
}

func (u *WalletUsecase) GetBalance(ctx context.Context, userID int64) (float64, error) {
	if _, err := u.userRepo.GetUserByID(ctx, userID); err != nil {
		return 0.0, err
	}
	return u.walletRepo.GetBalance(ctx, userID)
}

func (u *WalletUsecase) GetTransactions(ctx context.Context, userID int64) ([]*models2.Transaction, error) {
	if _, err := u.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	return u.walletRepo.GetTransactions(ctx, userID)
}
