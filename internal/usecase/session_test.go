package usecase_test

import (
	"computer-club/internal/repository/models"
	"computer-club/internal/usecase"
	"computer-club/mocks"
	"computer-club/pkg/errors"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartSession(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.UserRepository, *mocks.SessionRepository, *mocks.ComputerRepository, *mocks.TariffRepository, *mocks.WalletRepository, *mocks.Transaction)
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}
				session := &models.Session{ID: 1, UserID: userID, PCNumber: pcNumber, TariffID: tariffID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariffID).Return(session, nil)
				sessionRepo.On("CacheSession", mock.Anything, session).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "User Not Found",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				userRepo.On("GetUserByID", mock.Anything, userID).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
		},
		{
			name: "User Has Active Session",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(true, nil)
			},
			expectedError: errors.ErrSessionActive,
		},
		{
			name: "PC not found",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				user := &models.User{ID: userID}
				pcNumber := 101

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, errors.ErrComputerNotFound)
			},
			expectedError: errors.ErrComputerNotFound,
		},
		{
			name: "PC Busy",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(true, nil)
			},
			expectedError: errors.ErrPCBusy,
		},
		{
			name: "Tariff not found",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return((*models.Tariff)(nil), errors.ErrTariffNotFound)
			},
			expectedError: errors.ErrTariffNotFound,
		},
		{
			name: "Balance not found",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(0.0, errors.ErrCheckBalance)
			},
			expectedError: errors.ErrCheckBalance,
		},
		{
			name: "Insufficient Funds",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(50), nil)

			},
			expectedError: errors.ErrInsufficientFunds,
		},
		{
			name: "Withdraw",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(errors.ErrWithdraw)
			},
			expectedError: errors.ErrWithdraw,
		},
		{
			name: "Create Transaction Failed",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, errors.ErrCreateTransaction)
			},
			expectedError: errors.ErrCreateTransaction,
		},
		{
			name: "Create Session Failed",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariffID).Return((*models.Session)(nil), errors.ErrCreatedSession)
			},
			expectedError: errors.ErrCreatedSession,
		},
		{
			name: "Create Session Failed Update Computer",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariffID).Return((*models.Session)(nil), errors.ErrUpdateComputerStatus)
			},
			expectedError: errors.ErrUpdateComputerStatus,
		},
		{
			name: "Commit Transaction Failed",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}
				session := &models.Session{ID: 1, UserID: userID, PCNumber: pcNumber, TariffID: tariffID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(errors.ErrCommitData)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariffID).Return(session, nil)
			},
			expectedError: errors.ErrCommitData,
		},
		{
			name: "Transaction Commit Failure",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}
				session := &models.Session{ID: 1, UserID: userID, PCNumber: pcNumber, TariffID: tariffID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(errors.ErrCacheSession)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("IsComputerAvailable", mock.Anything, tx, pcNumber).Return(false, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariffID).Return(session, nil)
				sessionRepo.On("CacheSession", mock.Anything, session).Return(nil)
			},
			expectedError: errors.ErrCacheSession,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockSessionRepo := new(mocks.SessionRepository)
			mockComputerRepo := new(mocks.ComputerRepository)
			mockTariffRepo := new(mocks.TariffRepository)
			mockWalletRepo := new(mocks.WalletRepository)
			mockTransaction := new(mocks.Transaction)

			sessionUsecase := usecase.NewSessionUsecase(
				mockSessionRepo,
				mockUserRepo,
				mockComputerRepo,
				mockTariffRepo,
				mockWalletRepo,
			)

			tt.mockSetup(mockUserRepo, mockSessionRepo, mockComputerRepo, mockTariffRepo, mockWalletRepo, mockTransaction)

			_, err := sessionUsecase.StartSession(ctx, 1, 101, 5)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
