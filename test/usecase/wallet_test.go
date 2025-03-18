package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/usecase"
	"computer-club/mocks"
	"computer-club/pkg/errors"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestWithdraw(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.UserRepository, *mocks.WalletRepository, *mocks.TariffRepository, int64, float64)
		expectedError error
		userID        int64
		amount        float64
	}{
		{
			name:   "Success",
			userID: 42,
			amount: 100,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, nil, userID, amount).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Invalid Amount",
			userID: 3,
			amount: -50,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
			},
			expectedError: errors.ErrInvalidAmount,
		},
		{
			name:   "Wallet Not Found",
			userID: 42,
			amount: 100,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				walletRepo.On("GetBalance", mock.Anything, userID).Return(0.0, errors.ErrCheckBalance)
			},
			expectedError: errors.ErrCheckBalance,
		},
		{
			name:   "Insufficient Funds",
			userID: 3,
			amount: 300,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
			},
			expectedError: errors.ErrInsufficientFunds,
		},
		{
			name:   "Withdraw Failed",
			userID: 33,
			amount: 50,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, nil, userID, amount).Return(errors.ErrWithdraw)
			},
			expectedError: errors.ErrWithdraw,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)
			mockTariffRepo := new(mocks.TariffRepository)

			walletUsecase := usecase.NewWalletUsecase(
				mockWalletRepo,
				mockTariffRepo,
				mockUserRepo,
			)
			tt.mockSetup(mockUserRepo, mockWalletRepo, mockTariffRepo, tt.userID, tt.amount)

			err := walletUsecase.Withdraw(ctx, tt.userID, tt.amount)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestPutMoneyOnWallet(t *testing.T) {
	tests := []struct {
		name                string
		userID              int64
		amount              float64
		mockSetup           func(*mocks.UserRepository, *mocks.WalletRepository, *mocks.TariffRepository, *mocks.Transaction, int64, float64)
		expectedError       error
		expectedTransaction *models.Transaction
	}{
		{
			name:   "Success",
			userID: 101,
			amount: 500.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariff *mocks.TariffRepository, tx *mocks.Transaction, userID int64, amount float64) {
				transaction := &models.Transaction{ID: 1, UserID: userID, Amount: amount, TariffID: -1, Type: models.Add}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(1000.0, nil)
				walletRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				walletRepo.On("Deposit", mock.Anything, tx, userID, amount).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, amount, string(models.Add), mock.Anything).Return(transaction, nil)
			},
			expectedError:       nil,
			expectedTransaction: &models.Transaction{ID: 1, Type: models.Add},
		},
		{
			name:   "Invalid Amount",
			userID: 202,
			amount: -50.0, // Ошибка из-за отрицательной суммы
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariff *mocks.TariffRepository, tx *mocks.Transaction, userID int64, amount float64) {
			},
			expectedError: errors.ErrInvalidAmount,
		},
		{
			name:   "User Not Found",
			userID: 303,
			amount: 200.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariff *mocks.TariffRepository, tx *mocks.Transaction, userID int64, amount float64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
		},
		{
			name:   "Wallet Not Found",
			userID: 404,
			amount: 100.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariff *mocks.TariffRepository, tx *mocks.Transaction, userID int64, amount float64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(0.0, errors.ErrCheckBalance)
			},
			expectedError: errors.ErrCheckBalance,
		},
		{
			name:   "Deposit Failed",
			userID: 505,
			amount: 750.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariff *mocks.TariffRepository, tx *mocks.Transaction, userID int64, amount float64) {
				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(1200.0, nil)
				walletRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				walletRepo.On("Deposit", mock.Anything, tx, userID, amount).Return(errors.ErrToDeposit)
			},
			expectedError: errors.ErrToDeposit,
		},
		{
			name:   "Commit Transaction Failed",
			userID: 606,
			amount: 300.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariff *mocks.TariffRepository, tx *mocks.Transaction, userID int64, amount float64) {
				transaction := &models.Transaction{ID: 1, UserID: userID, Amount: amount, TariffID: -1, Type: models.Add}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(errors.ErrCommitData)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(800.0, nil)
				walletRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				walletRepo.On("Deposit", mock.Anything, tx, userID, amount).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, amount, string(models.Add), mock.Anything).Return(transaction, nil)
			},
			expectedError: errors.ErrCommitData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)
			mockTariffRepo := new(mocks.TariffRepository)
			mockTransaction := new(mocks.Transaction)

			walletUsecase := usecase.NewWalletUsecase(
				mockWalletRepo,
				mockTariffRepo,
				mockUserRepo,
			)

			tt.mockSetup(mockUserRepo, mockWalletRepo, mockTariffRepo, mockTransaction, tt.userID, tt.amount)

			transaction, err := walletUsecase.PutMoneyOnWallet(ctx, tt.userID, tt.amount)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, transaction.UserID)
				assert.Equal(t, tt.amount, transaction.Amount)
				assert.Equal(t, string(models.Add), string(transaction.Type))
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		mockSetup     func(*mocks.UserRepository, *mocks.WalletRepository, int64)
		expectedError error
		expectedData  []*models.Transaction
	}{
		{
			name:   "Success",
			userID: 42,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, userID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				transactions := []*models.Transaction{
					{ID: 1, UserID: userID, Amount: 100.00, TariffID: -1, Type: models.Buy},
					{ID: 2, UserID: userID, Amount: 50.00, TariffID: -1, Type: models.Add},
				}
				walletRepo.On("GetTransactions", mock.Anything, userID).Return(transactions, nil)
			},
			expectedError: nil,
			expectedData: []*models.Transaction{
				{ID: 1, UserID: 42, Amount: 100.00, TariffID: -1, Type: models.Buy},
				{ID: 2, UserID: 42, Amount: 50.00, TariffID: -1, Type: models.Add},
			},
		},
		{
			name:   "User Not Found",
			userID: 99,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, userID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
			expectedData:  nil,
		},
		{
			name:   "No Transactions Found",
			userID: 55,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, userID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetTransactions", mock.Anything, userID).Return(nil, nil)
			},
			expectedError: nil,
			expectedData:  nil,
		},
		{
			name:   "Wallet Repository Error",
			userID: 88,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, userID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetTransactions", mock.Anything, userID).Return(nil, errors.ErrCheckTransaction)
			},
			expectedError: errors.ErrCheckTransaction,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)

			walletUsecase := usecase.NewWalletUsecase(
				mockWalletRepo,
				nil,
				mockUserRepo,
			)

			tt.mockSetup(mockUserRepo, mockWalletRepo, tt.userID)

			result, err := walletUsecase.GetTransactions(ctx, tt.userID)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedData, result)
		})
	}
}

func TestGetBalance(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		mockSetup     func(*mocks.UserRepository, *mocks.WalletRepository, int64)
		expectedError error
		expectedData  float64
	}{
		{
			name:   "Success",
			userID: 42,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, userID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(500.75, nil)
			},
			expectedError: nil,
			expectedData:  500.75,
		},
		{
			name:   "User Not Found",
			userID: 99,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, userID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
			expectedData:  0.0,
		},
		{
			name:   "Wallet Repository Error",
			userID: 88,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, userID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(0.0, errors.ErrCheckBalance)
			},
			expectedError: errors.ErrCheckBalance,
			expectedData:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)

			walletUsecase := usecase.NewWalletUsecase(
				mockWalletRepo,
				nil,
				mockUserRepo,
			)

			tt.mockSetup(mockUserRepo, mockWalletRepo, tt.userID)

			result, err := walletUsecase.GetBalance(ctx, tt.userID)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedData, result)
		})
	}
}
