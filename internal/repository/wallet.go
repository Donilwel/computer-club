package repository

import (
	"computer-club/internal/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepository interface {
	GetBalance(ctx context.Context, userID int64) (float64, error)
	GetTransactions(ctx context.Context, userID int64) ([]*models.Transaction, error)
	CreateWallet(ctx context.Context, tx Transaction, wallet *models.Wallet) error
	Deposit(ctx context.Context, tx Transaction, userID int64, amount float64) error
	Withdraw(ctx context.Context, tx Transaction, userID int64, amount float64) error
	CreateTransaction(ctx context.Context, tx Transaction, userID int64,
		amount float64, typ string, tariff *models.Tariff) (*models.Transaction, error)
	BeginTransaction(ctx context.Context) Transaction
}

type PostgresWalletRepo struct {
	db *gorm.DB
}

func NewPostgresWalletRepo(db *gorm.DB) WalletRepository {
	return &PostgresWalletRepo{db: db}
}

func (r *PostgresWalletRepo) BeginTransaction(ctx context.Context) Transaction {
	return &GormTransaction{tx: r.db.WithContext(ctx).Begin()}
}

func (r *PostgresWalletRepo) CreateTransaction(ctx context.Context, tx Transaction, userID int64,
	amount float64, typ string, tariff *models.Tariff) (*models.Transaction, error) {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	var tariffID int64
	if tariff == nil {
		tariffID = -1
	} else {
		tariffID = tariff.ID
	}

	transaction := models.Transaction{
		UserID:   userID,
		Amount:   amount,
		Type:     models.TransactionType(typ),
		TariffID: tariffID,
	}
	if err := db.WithContext(ctx).Create(&transaction).Error; err != nil {
		return nil, errors.ErrCreateTransaction
	}
	return &transaction, nil
}

func (r *PostgresWalletRepo) Deposit(ctx context.Context, tx Transaction, userID int64, amount float64) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	var wallet models.Wallet
	err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&wallet).Error
	if err != nil {
		return errors.ErrUnexpected
	}

	err = db.WithContext(ctx).
		Model(&wallet).
		Where("user_id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
	if err != nil {
		return errors.ErrToDeposit
	}
	return nil
}

func (r *PostgresWalletRepo) Withdraw(ctx context.Context, tx Transaction, userID int64, amount float64) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	var wallet models.Wallet
	err := db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		First(&wallet).Error
	if err != nil {
		return errors.ErrUnexpected
	}
	err = db.WithContext(ctx).Model(&wallet).
		Where("user_id = ? AND balance >= ?", userID, amount).
		Update("balance", gorm.Expr("balance - ?", amount)).Error
	if err != nil {
		return errors.ErrWithdraw
	}
	return nil
}

func (r *PostgresWalletRepo) GetBalance(ctx context.Context, userID int64) (float64, error) {
	var wallet models.Wallet
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return 0.0, errors.ErrCheckBalance
	}
	return wallet.Balance, nil
}

func (r *PostgresWalletRepo) GetTransactions(ctx context.Context, userID int64) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
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

func (r *PostgresWalletRepo) CreateWallet(ctx context.Context, tx Transaction, wallet *models.Wallet) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	if err := db.WithContext(ctx).Create(wallet).Error; err != nil {
		return errors.ErrCreateWallet
	}
	return nil
}
