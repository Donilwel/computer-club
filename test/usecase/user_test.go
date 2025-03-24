package usecase_test

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

func TestGetInfoUser(t *testing.T) {
	tests := []struct {
		name                 string
		mockSetup            func(*mocks.UserRepository, *mocks.WalletRepository)
		expectedError        error
		expectedUser         *models.User
		expectedBalance      float64
		expectedTransactions []*models.Transaction
	}{
		{
			name: "Success",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				userID := int64(1)
				user := &models.User{ID: userID, Name: "Alice", Email: "alice@example.com"}
				balance := 500.0
				transactions := []*models.Transaction{
					{ID: 1, UserID: userID, Amount: 200, Type: models.Buy},
					{ID: 2, UserID: userID, Amount: 300, Type: models.Add},
				}

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(balance, nil)
				walletRepo.On("GetTransactions", mock.Anything, userID).Return(transactions, nil)
			},
			expectedError:   nil,
			expectedUser:    &models.User{ID: 1, Name: "Alice", Email: "alice@example.com"},
			expectedBalance: 500.0,
			expectedTransactions: []*models.Transaction{
				{ID: 1, UserID: 1, Amount: 200, Type: models.Buy},
				{ID: 2, UserID: 1, Amount: 300, Type: models.Add},
			},
		},
		{
			name: "User Not Found",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				userID := int64(1)
				userRepo.On("GetUserByID", mock.Anything, userID).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
		},
		{
			name: "Transaction Fetch Failed",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				userID := int64(1)
				user := &models.User{ID: userID, Name: "Alice", Email: "alice@example.com"}
				balance := 500.0

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(balance, nil)
				walletRepo.On("GetTransactions", mock.Anything, userID).Return([]*models.Transaction{}, errors.ErrCheckTransaction)
			},
			expectedError: errors.ErrCheckTransaction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)

			userUsecase := usecase.NewUserUsecase(
				mockUserRepo,
				mockWalletRepo,
			)

			tt.mockSetup(mockUserRepo, mockWalletRepo)

			user, balance, transactions, err := userUsecase.GetInfoUser(ctx, 1)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
				assert.Equal(t, tt.expectedBalance, balance)
				assert.Equal(t, tt.expectedTransactions, transactions)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		mockSetup     func(*mocks.UserRepository)
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:  "Success",
			email: "johndoe@example.com",
			mockSetup: func(userRepo *mocks.UserRepository) {
				user := &models.User{ID: 1, Email: "johndoe@example.com"}
				userRepo.On("GetUserByEmail", mock.Anything, "johndoe@example.com").Return(user, nil)
			},
			expectedUser:  &models.User{ID: 1, Email: "johndoe@example.com"},
			expectedError: nil,
		},
		{
			name:  "User Not Found",
			email: "notfound@example.com",
			mockSetup: func(userRepo *mocks.UserRepository) {
				userRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedUser:  nil,
			expectedError: errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			walletRepo := new(mocks.WalletRepository)

			userUsecase := usecase.NewUserUsecase(mockUserRepo, walletRepo)

			tt.mockSetup(mockUserRepo)

			user, err := userUsecase.GetUserByEmail(ctx, tt.email)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedUser, user)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		mockSetup     func(*mocks.UserRepository)
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:   "Success",
			userID: 1,
			mockSetup: func(userRepo *mocks.UserRepository) {
				user := &models.User{ID: 1, Email: "johndoe@example.com"}
				userRepo.On("GetUserByID", mock.Anything, int64(1)).Return(user, nil)
			},
			expectedUser:  &models.User{ID: 1, Email: "johndoe@example.com"},
			expectedError: nil,
		},
		{
			name:   "User Not Found",
			userID: 2,
			mockSetup: func(userRepo *mocks.UserRepository) {
				userRepo.On("GetUserByID", mock.Anything, int64(2)).Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedUser:  nil,
			expectedError: errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			walletRepo := new(mocks.WalletRepository)

			userUsecase := usecase.NewUserUsecase(mockUserRepo, walletRepo)

			tt.mockSetup(mockUserRepo)

			user, err := userUsecase.GetUserByID(ctx, tt.userID)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedUser, user)
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.UserRepository)
		expectedUsers []*models.User
		expectedError error
	}{
		{
			name: "Success",
			mockSetup: func(userRepo *mocks.UserRepository) {
				users := []*models.User{
					{ID: 1, Name: "John Doe", Role: "admin", Email: "john@example.com"},
					{ID: 2, Name: "Jane Doe", Role: "customer", Email: "jane@example.com"},
				}
				userRepo.On("GetUsers", mock.Anything).Return(users, nil)
			},
			expectedUsers: []*models.User{
				{ID: 1, Name: "John Doe", Role: "admin", Email: "john@example.com"},
				{ID: 2, Name: "Jane Doe", Role: "customer", Email: "jane@example.com"},
			},
			expectedError: nil,
		},
		{
			name: "Database Error",
			mockSetup: func(userRepo *mocks.UserRepository) {
				userRepo.On("GetUsers", mock.Anything).Return(nil, errors.ErrFindUsers)
			},
			expectedUsers: nil,
			expectedError: errors.ErrFindUsers,
		},
		{
			name: "No Users Found",
			mockSetup: func(userRepo *mocks.UserRepository) {
				userRepo.On("GetUsers", mock.Anything).Return(nil, errors.ErrZeroUsers)
			},
			expectedUsers: nil,
			expectedError: errors.ErrZeroUsers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)
			tt.mockSetup(mockUserRepo)

			userUsecase := usecase.NewUserUsecase(mockUserRepo, mockWalletRepo)
			users, err := userUsecase.GetUsers(ctx)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedUsers, users)
		})
	}
}
