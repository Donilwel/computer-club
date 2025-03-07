package usecase

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"computer-club/internal/repository"
)

type WalletService interface {
	Deposit(userID int64, amount float64) error
	Withdraw(userID int64, amount float64) error
	GetBalance(userID int64) (float64, error)
	GetTransactions(userID int64) ([]models.Transaction, error)
	ChargeForSession(userID int64, tariffID int64) error
	CreateWallet(userID int64) error
}

type WalletUsecase struct {
	walletRepo repository.WalletRepository
	tariffRepo repository.TariffRepository
}

func (u *WalletUsecase) CreateWallet(userID int64) error {
	_, err := u.walletRepo.GetBalance(userID)
	if err == nil {
		return errors.ErrWalletAlreadyExists
	}

	wallet := &models.Wallet{
		UserID:  userID,
		Balance: 0.0,
	}

	return u.walletRepo.CreateWallet(wallet)
}

func (u *WalletUsecase) ChargeForSession(userID int64, tariffID int64) error {
	tariff, err := u.tariffRepo.GetTariffByID(tariffID)
	if err != nil {
		return errors.ErrTariffNotFound
	}

	// Проверяем баланс пользователя
	balance, err := u.walletRepo.GetBalance(userID)
	if err != nil {
		return err
	}
	if balance < tariff.Price {
		return errors.ErrInsufficientFunds
	}

	// Списываем средства
	if err := u.walletRepo.Withdraw(userID, tariff.Price); err != nil {
		return err
	}

	return nil
}

func NewWalletUsecase(walletRepo repository.WalletRepository, tariffRepo repository.TariffRepository) *WalletUsecase {
	return &WalletUsecase{walletRepo: walletRepo, tariffRepo: tariffRepo}
}

func (u *WalletUsecase) Deposit(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.ErrInvalidAmount
	}
	return u.walletRepo.Deposit(userID, amount)
}

func (u *WalletUsecase) Withdraw(userID int64, amount float64) error {
	if amount <= 0 {
		return errors.ErrInvalidAmount
	}
	balance, err := u.walletRepo.GetBalance(userID)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.ErrInsufficientFunds
	}
	return u.walletRepo.Withdraw(userID, amount)
}

func (u *WalletUsecase) GetBalance(userID int64) (float64, error) {
	return u.walletRepo.GetBalance(userID)
}

func (u *WalletUsecase) GetTransactions(userID int64) ([]models.Transaction, error) {
	return u.walletRepo.GetTransactions(userID)
}
