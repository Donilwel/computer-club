package usecase

import (
	"computer-club/internal/repository"
	models2 "computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
)

type WalletService interface {
	Deposit(ctx context.Context, userID int64, amount float64) error
	Withdraw(ctx context.Context, userID int64, amount float64) error
	GetBalance(ctx context.Context, userID int64) (float64, error)
	GetTransactions(ctx context.Context, userID int64) ([]models2.Transaction, error)
	CreateTransaction(ctx context.Context,
		userID int64,
		amount float64,
		typ string,
		tariffID int64) (*models2.Transaction, error)
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

func (u *WalletUsecase) Deposit(ctx context.Context, userID int64, amount float64) error {
	if amount <= 0 {
		return errors.ErrInvalidAmount
	}

	if _, err := u.userRepo.GetUserByID(ctx, userID); err != nil {
		return err
	}
	if _, err := u.walletRepo.GetBalance(ctx, userID); err != nil {
		return err
	}

	return u.walletRepo.Deposit(ctx, userID, amount)
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

func (u *WalletUsecase) GetTransactions(ctx context.Context, userID int64) ([]models2.Transaction, error) {
	if _, err := u.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	return u.walletRepo.GetTransactions(ctx, userID)
}

func (u *WalletUsecase) CreateTransaction(ctx context.Context,
	userID int64,
	amount float64,
	typ string,
	tariffID int64) (*models2.Transaction, error) {
	if _, err := u.userRepo.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, errors.ErrInvalidAmount
	}
	if typ != string(models2.Buy) && typ != string(models2.Add) {
		return nil, errors.ErrorTypeTransaction
	}
	if tariffID != -1 {
		tariff, err := u.tariffRepo.GetTariffByID(ctx, tariffID)
		if err != nil {
			return nil, err
		}
		return u.walletRepo.CreateTransaction(ctx, nil, userID, amount, typ, tariff)
	}
	return u.walletRepo.CreateTransaction(ctx, nil, userID, amount, typ, nil)
}
