package repository

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Deposit(userID int64, amount float64) error
	Withdraw(userID int64, amount float64) error
	GetBalance(userID int64) (float64, error)
	GetTransactions(userID int64) ([]models.Transaction, error)
	CreateWallet(wallet *models.Wallet) error
}

type PostgresWalletRepo struct {
	db *gorm.DB
}

func NewPostgresWalletRepo(db *gorm.DB) *PostgresWalletRepo {
	return &PostgresWalletRepo{db: db}
}

func (r *PostgresWalletRepo) Deposit(userID int64, amount float64) error {
	err := r.db.Model(&models.Wallet{}).Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return errors.ErrToDeposit
	}
	return nil
}

func (r *PostgresWalletRepo) Withdraw(userID int64, amount float64) error {
	err := r.db.Model(&models.Wallet{}).Where("user_id = ? AND balance >= ?", userID, amount).
		Update("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		return errors.ErrWithdraw
	}
	return nil
}

func (r *PostgresWalletRepo) GetBalance(userID int64) (float64, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return 0, errors.ErrCheckBalance
	}
	return wallet.Balance, nil
}

func (r *PostgresWalletRepo) GetTransactions(userID int64) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("user_id = ?", userID).Find(&transactions).Error
	if err != nil {
		return nil, errors.ErrCheckTransaction
	}
	return transactions, nil
}

func (r *PostgresWalletRepo) CreateWallet(wallet *models.Wallet) error {
	if err := r.db.Create(wallet).Error; err != nil {
		return errors.ErrCreateWallet
	}
	return nil
}
