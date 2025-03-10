package repository

import (
	models2 "computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetBalance(ctx context.Context, userID int64) (float64, error)
	GetTransactions(ctx context.Context, userID int64) ([]models2.Transaction, error)
	CreateWallet(ctx context.Context, wallet *models2.Wallet) error
	Deposit(ctx context.Context, userID int64, amount float64) error
	Withdraw(tx *gorm.DB, userID int64, amount float64) error
	CreateTransaction(tx *gorm.DB, userID int64, amount float64, typ string, tariff *models2.Tariff) (*models2.Transaction, error)
}

type PostgresWalletRepo struct {
	db *gorm.DB
}

func NewPostgresWalletRepo(db *gorm.DB) WalletRepository {
	return &PostgresWalletRepo{db: db}
}

func (r *PostgresWalletRepo) CreateTransaction(tx *gorm.DB, userID int64, amount float64, typ string, tariff *models2.Tariff) (*models2.Transaction, error) {
	if tx == nil {
		tx = r.db
	}
	var tariffID int64
	if tariff == nil {
		tariffID = -1
	} else {
		tariffID = tariff.ID
	}

	transaction := models2.Transaction{
		UserID:   userID,
		Amount:   amount,
		Type:     models2.TransactionType(typ),
		TariffID: tariffID,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		return nil, errors.ErrCreateTransaction
	}
	return &transaction, nil
}

func (r *PostgresWalletRepo) Deposit(ctx context.Context, userID int64, amount float64) error {
	err := r.db.WithContext(ctx).
		Model(&models2.Wallet{}).
		Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return errors.ErrToDeposit
	}
	return nil
}

func (r *PostgresWalletRepo) Withdraw(tx *gorm.DB, userID int64, amount float64) error {
	if tx == nil {
		tx = r.db
	}
	err := tx.Model(&models2.Wallet{}).
		Where("user_id = ? AND balance >= ?", userID, amount).
		Update("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		return errors.ErrWithdraw
	}
	return nil
}

func (r *PostgresWalletRepo) GetBalance(ctx context.Context, userID int64) (float64, error) {
	var wallet models2.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return 0, errors.ErrCheckBalance
	}
	return wallet.Balance, nil
}

func (r *PostgresWalletRepo) GetTransactions(ctx context.Context, userID int64) ([]models2.Transaction, error) {
	var transactions []models2.Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(50).
		Find(&transactions).Error
	if err != nil {
		return nil, errors.ErrCheckTransaction
	}
	return transactions, nil
}

func (r *PostgresWalletRepo) CreateWallet(ctx context.Context, wallet *models2.Wallet) error {
	if err := r.db.WithContext(ctx).Create(wallet).Error; err != nil {
		return errors.ErrCreateWallet
	}
	return nil
}
