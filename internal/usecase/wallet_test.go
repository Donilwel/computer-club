package usecase

import (
	"computer-club/internal/repository/models"
	"computer-club/mocks"
	"computer-club/pkg/errors"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestDeposit(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.UserRepository, *mocks.WalletRepository, *mocks.TariffRepository, int64, float64)
		expectedError error
		userID        int64
		amount        float64
	}{
		{
			name:   "Success",
			userID: 42, amount: 150.50,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(300), nil)
				walletRepo.On("Deposit", mock.Anything, userID, amount).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Invalid Amount",
			userID: 10, amount: -50.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
			},
			expectedError: errors.ErrInvalidAmount,
		},
		{
			name:   "User Not Found",
			userID: 99, amount: 200.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
		},
		{
			name:   "Wallet Not Found",
			userID: 13, amount: 500.0,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(0), errors.ErrCheckBalance)
			},
			expectedError: errors.ErrCheckBalance,
		},
		{
			name:   "Deposit Failed",
			userID: 77, amount: 75.25,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Deposit", mock.Anything, userID, amount).Return(errors.ErrToDeposit)
			},
			expectedError: errors.ErrToDeposit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)
			mockTariffRepo := new(mocks.TariffRepository)

			walletUsecase := NewWalletUsecase(
				mockWalletRepo,
				mockTariffRepo,
				mockUserRepo,
			)

			tt.mockSetup(mockUserRepo, mockWalletRepo, mockTariffRepo, tt.userID, tt.amount)

			err := walletUsecase.Deposit(ctx, tt.userID, tt.amount)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

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

			walletUsecase := NewWalletUsecase(
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

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		amount        float64
		typ           string
		tariffID      int64
		mockSetup     func(*mocks.UserRepository, *mocks.WalletRepository, *mocks.TariffRepository, int64, float64, string, int64)
		expectedError error
	}{
		{
			name:     "Success with Tariff",
			userID:   42,
			amount:   150.75,
			typ:      string(models.Buy),
			tariffID: 5,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64, typ string, tariffID int64) {
				tariff := &models.Tariff{ID: tariffID, Price: amount}
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("CreateTransaction", mock.Anything, nil, userID, amount, typ, tariff).Return(&models.Transaction{ID: 1}, nil)
			},
			expectedError: nil,
		},
		{
			name:     "Success without Tariff",
			userID:   55,
			amount:   200.00,
			typ:      string(models.Add),
			tariffID: -1,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64, typ string, tariffID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("CreateTransaction", mock.Anything, nil, userID, amount, typ, mock.MatchedBy(func(t *models.Tariff) bool { return t == nil })).Return(&models.Transaction{ID: 2}, nil)
			},
			expectedError: nil,
		},
		{
			name:     "User Not Found",
			userID:   99,
			amount:   100.00,
			typ:      string(models.Buy),
			tariffID: 5,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64, typ string, tariffID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
		},
		{
			name:     "Invalid Amount",
			userID:   88,
			amount:   -10,
			typ:      string(models.Buy),
			tariffID: -1,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64, typ string, tariffID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
			},
			expectedError: errors.ErrInvalidAmount,
		},
		{
			name:     "Invalid Transaction Type",
			userID:   77,
			amount:   50,
			typ:      "invalid",
			tariffID: -1,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64, typ string, tariffID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
			},
			expectedError: errors.ErrorTypeTransaction,
		},
		{
			name:     "Tariff not found",
			userID:   88,
			amount:   100.00,
			typ:      string(models.Buy),
			tariffID: 10,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64, typ string, tariffID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return((*models.Tariff)(nil), errors.ErrTariffNotFound)
			},
			expectedError: errors.ErrTariffNotFound,
		},
		{
			name:     "Failed to Create Transaction",
			userID:   99,
			amount:   500.00,
			typ:      string(models.Buy),
			tariffID: -1,
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tariffRepo *mocks.TariffRepository, userID int64, amount float64, typ string, tariffID int64) {
				userRepo.On("GetUserByID", mock.Anything, userID).Return(&models.User{ID: userID}, nil)
				walletRepo.On("CreateTransaction", mock.Anything, nil, userID, amount, typ, mock.MatchedBy(func(t *models.Tariff) bool { return t == nil })).Return(nil, errors.ErrCreateTransaction)
			},
			expectedError: errors.ErrCreateTransaction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)
			mockTariffRepo := new(mocks.TariffRepository)

			walletUsecase := NewWalletUsecase(
				mockWalletRepo,
				mockTariffRepo,
				mockUserRepo,
			)

			tt.mockSetup(mockUserRepo, mockWalletRepo, mockTariffRepo, tt.userID, tt.amount, tt.typ, tt.tariffID)

			_, err := walletUsecase.CreateTransaction(ctx, tt.userID, tt.amount, tt.typ, tt.tariffID)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
