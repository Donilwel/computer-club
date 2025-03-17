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
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}
				session := &models.Session{ID: 1, UserID: userID, PCNumber: pcNumber, TariffID: tariffID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariff).Return(session, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Busy)).Return(nil)
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
				pcNumber := 101
				user := &models.User{ID: userID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(nil, errors.ErrComputerNotFound)
			},
			expectedError: errors.ErrComputerNotFound,
		},
		{
			name: "PC Busy",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				user := &models.User{ID: userID}
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Busy}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
			},
			expectedError: errors.ErrPCBusy,
		},
		{
			name: "Tariff not found",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, errors.ErrTariffNotFound)
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
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
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
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
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
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
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
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
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
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariff).Return(nil, errors.ErrCreatedSession)
			},
			expectedError: errors.ErrCreatedSession,
		},
		{
			name: "Failed Update Computer status",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}
				session := &models.Session{ID: 1, UserID: userID, PCNumber: pcNumber, TariffID: tariffID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariff).Return(session, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Busy)).Return(errors.ErrUpdateComputerStatus)
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
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}
				session := &models.Session{ID: 1, UserID: userID, PCNumber: pcNumber, TariffID: tariffID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(errors.ErrCommitData)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariff).Return(session, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Busy)).Return(nil)
				sessionRepo.On("CacheSession", mock.Anything, session).Return(nil)
			},
			expectedError: errors.ErrCommitData,
		},
		{
			name: "Cashed Session Failure",
			mockSetup: func(userRepo *mocks.UserRepository, sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tariffRepo *mocks.TariffRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userID := int64(1)
				pcNumber := 101
				tariffID := int64(5)
				tariff := &models.Tariff{ID: tariffID, Price: 100}
				user := &models.User{ID: userID}
				computer := &models.Computer{ID: int64(pcNumber), PCNumber: pcNumber, Status: models.Free}
				session := &models.Session{ID: 1, UserID: userID, PCNumber: pcNumber, TariffID: tariffID}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("HasActiveSession", mock.Anything, tx, userID).Return(false, nil)
				computerRepo.On("GetComputerByID", mock.Anything, pcNumber).Return(computer, nil)
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(float64(200), nil)
				walletRepo.On("Withdraw", mock.Anything, tx, userID, tariff.Price).Return(nil)
				walletRepo.On("CreateTransaction", mock.Anything, tx, userID, tariff.Price, mock.Anything, tariff).Return(nil, nil)
				sessionRepo.On("CreateSession", mock.Anything, tx, userID, pcNumber, tariff).Return(session, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Busy)).Return(nil)
				sessionRepo.On("CacheSession", mock.Anything, session).Return(errors.ErrCacheSession)
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

func TestEndSession(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.SessionRepository, *mocks.ComputerRepository, *mocks.Transaction)
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Active}
				computer := &models.Computer{ID: int64(session.PCNumber), PCNumber: session.PCNumber, Status: models.Busy}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
				sessionRepo.On("MarkSessionFinished", mock.Anything, tx, sessionID).Return(nil)
				computerRepo.On("GetComputerByID", mock.Anything, session.PCNumber).Return(computer, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Free)).Return(nil)
				sessionRepo.On("DeleteSessionCache", mock.Anything, sessionID).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Session Not Found",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)

				tx.On("Rollback").Return(nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(nil, errors.ErrSessionNotFound)
			},
			expectedError: errors.ErrSessionNotFound,
		},
		{
			name: "Session Already Finished",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Finished}

				tx.On("Rollback").Return(nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
			},
			expectedError: errors.ErrStatusSessionAlreadyFinished,
		},
		{
			name: "Upload Status Session Failed",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Active}

				tx.On("Rollback").Return(nil)
				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
				sessionRepo.On("MarkSessionFinished", mock.Anything, tx, sessionID).Return(errors.ErrUpdateSession)
			},
			expectedError: errors.ErrUpdateSession,
		},
		{
			name: "Computer Not Found",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Active}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
				sessionRepo.On("MarkSessionFinished", mock.Anything, tx, sessionID).Return(nil)
				computerRepo.On("GetComputerByID", mock.Anything, session.PCNumber).Return(nil, errors.ErrComputerNotFound)
			},
			expectedError: errors.ErrComputerNotFound,
		},
		{
			name: "Status Computer Already Free",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Active}
				computer := &models.Computer{ID: int64(session.PCNumber), PCNumber: session.PCNumber, Status: models.Free}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
				sessionRepo.On("MarkSessionFinished", mock.Anything, tx, sessionID).Return(nil)
				computerRepo.On("GetComputerByID", mock.Anything, session.PCNumber).Return(computer, nil)
			},
			expectedError: errors.ErrComputerAlreadyFree,
		},
		{
			name: "Update Computer Status Failed",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Active}
				computer := &models.Computer{ID: int64(session.PCNumber), PCNumber: session.PCNumber, Status: models.Busy}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
				sessionRepo.On("MarkSessionFinished", mock.Anything, tx, sessionID).Return(nil)
				computerRepo.On("GetComputerByID", mock.Anything, session.PCNumber).Return(computer, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Free)).Return(errors.ErrUpdateSession)
			},
			expectedError: errors.ErrUpdateSession,
		},
		{
			name: "Transaction Commit Failed",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Active}
				computer := &models.Computer{ID: int64(session.PCNumber), PCNumber: session.PCNumber, Status: models.Busy}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(errors.ErrCommitData)

				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
				sessionRepo.On("MarkSessionFinished", mock.Anything, tx, sessionID).Return(nil)
				computerRepo.On("GetComputerByID", mock.Anything, session.PCNumber).Return(computer, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Free)).Return(nil)
			},
			expectedError: errors.ErrCommitData,
		},
		{
			name: "Delete Cashed Data Failed",
			mockSetup: func(sessionRepo *mocks.SessionRepository, computerRepo *mocks.ComputerRepository, tx *mocks.Transaction) {
				sessionID := int64(1)
				session := &models.Session{ID: sessionID, PCNumber: 5, Status: models.Active}
				computer := &models.Computer{ID: int64(session.PCNumber), PCNumber: session.PCNumber, Status: models.Busy}

				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				sessionRepo.On("BeginTransaction", mock.Anything).Return(tx, nil)
				sessionRepo.On("GetSessionByID", mock.Anything, tx, sessionID).Return(session, nil)
				sessionRepo.On("MarkSessionFinished", mock.Anything, tx, sessionID).Return(nil)
				computerRepo.On("GetComputerByID", mock.Anything, session.PCNumber).Return(computer, nil)
				computerRepo.On("ChangeComputerStatus", mock.Anything, tx, computer, string(models.Free)).Return(nil)
				sessionRepo.On("DeleteSessionCache", mock.Anything, sessionID).Return(errors.ErrDeleteCashedData)
			},
			expectedError: errors.ErrDeleteCashedData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockSessionRepo := new(mocks.SessionRepository)
			mockComputerRepo := new(mocks.ComputerRepository)
			mockTransaction := new(mocks.Transaction)

			sessionUsecase := usecase.NewSessionUsecase(
				mockSessionRepo,
				new(mocks.UserRepository),
				mockComputerRepo,
				new(mocks.TariffRepository),
				new(mocks.WalletRepository),
			)

			tt.mockSetup(mockSessionRepo, mockComputerRepo, mockTransaction)

			err := sessionUsecase.EndSession(ctx, 1)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestGetActiveSessions(t *testing.T) {
	mockSessionRepo := new(mocks.SessionRepository)
	sessionUsecase := usecase.NewSessionUsecase(
		mockSessionRepo,
		new(mocks.UserRepository),
		new(mocks.ComputerRepository),
		new(mocks.TariffRepository),
		new(mocks.WalletRepository),
	)

	ctx := context.Background()
	activeSessions := []*models.Session{
		{ID: 1, UserID: 101, PCNumber: 5, TariffID: 2},
		{ID: 2, UserID: 102, PCNumber: 6, TariffID: 3},
	}

	mockSessionRepo.On("GetActiveSessions", ctx).Return(activeSessions)

	result := sessionUsecase.GetActiveSessions(ctx)

	assert.Equal(t, activeSessions, result)

	mockSessionRepo.AssertExpectations(t)
}
