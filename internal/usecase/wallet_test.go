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
