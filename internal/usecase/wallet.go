package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"computer-club/pkg/errors"
	"context"
)

type WalletService interface {
	Deposit(ctx context.Context, userID int64, amount float64) error
	Withdraw(ctx context.Context, userID int64, amount float64) error
	GetBalance(ctx context.Context, userID int64) (float64, error)
	GetTransactions(ctx context.Context, userID int64) ([]models.Transaction, error)
	ChargeForSession(ctx context.Context, userID int64, tariffID int64) error
	CreateTransaction(ctx context.Context, userID int64, amount float64, typ string, tariffID int64) (*models.Transaction, error)
	CreateWallet(ctx context.Context, userID int64) error
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

func (u *WalletUsecase) CreateWallet(ctx context.Context, userID int64) error {
	_, err := u.walletRepo.GetBalance(ctx, userID)
	if err == nil {
		return errors.ErrWalletAlreadyExists
	}

	wallet := &models.Wallet{
		UserID:  userID,
		Balance: 0.0,
	}

	return u.walletRepo.CreateWallet(ctx, wallet)
}

func (u *WalletUsecase) ChargeForSession(ctx context.Context, userID int64, tariffID int64) error {
	tariff, err := u.tariffRepo.GetTariffByID(ctx, tariffID)
	if err != nil {
		return errors.ErrTariffNotFound
	}

	// Проверяем баланс пользователя
	balance, err := u.walletRepo.GetBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance < tariff.Price {
		return errors.ErrInsufficientFunds
	}
	// Списываем средства
	if err := u.walletRepo.Withdraw(ctx, userID, tariff.Price); err != nil {
		return err
	}

	_, err = u.walletRepo.CreateTransaction(ctx, userID, tariff.Price, string(models.Buy), tariff)
	if err != nil {
		return err
	}

	return nil
}

func (u *WalletUsecase) Deposit(ctx context.Context, userID int64, amount float64) error {
	if amount <= 0 {
		return errors.ErrInvalidAmount
	}

	if _, err := u.userRepo.GetUserByID(ctx, userID); err != nil {
		return errors.ErrUserNotFound
	}
	if _, err := u.walletRepo.GetBalance(ctx, userID); err != nil {
		return errors.ErrCheckBalance
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
	return u.walletRepo.Withdraw(ctx, userID, amount)
}

func (u *WalletUsecase) GetBalance(ctx context.Context, userID int64) (float64, error) {
	return u.walletRepo.GetBalance(ctx, userID)
}

func (u *WalletUsecase) GetTransactions(ctx context.Context, userID int64) ([]models.Transaction, error) {
	return u.walletRepo.GetTransactions(ctx, userID)
}

func (u *WalletUsecase) CreateTransaction(ctx context.Context, userID int64, amount float64, typ string, tariffID int64) (*models.Transaction, error) {
	if tariffID != -1 {
		tariff, err := u.tariffRepo.GetTariffByID(ctx, tariffID)
		if err != nil {
			return nil, err
		}
		return u.walletRepo.CreateTransaction(ctx, userID, amount, typ, tariff)
	}
	return u.walletRepo.CreateTransaction(ctx, userID, amount, typ, nil)
}
