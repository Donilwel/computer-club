package usecase_test

import (
	"bou.ke/monkey"
	"computer-club/internal/repository/models"
	"computer-club/internal/usecase"
	"computer-club/mocks"
	"computer-club/pkg/JWT"
	"computer-club/pkg/errors"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			name     string
			email    string
			password string
			role     models.UserRole
		}
		mockSetup     func(*mocks.UserRepository, *mocks.WalletRepository, *mocks.Transaction)
		expectedError error
	}{
		{
			name: "Success",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "JohnDoe",
				email:    "johndoe@example.com",
				password: "securepass",
				role:     models.Customer,
			},
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("BeginTransaction", mock.Anything).Return(tx)
				userRepo.On("GetUserByEmail", mock.Anything, "johndoe@example.com").Return((*models.User)(nil), nil)
				userRepo.On("GetUserByName", mock.Anything, "JohnDoe").Return((*models.User)(nil), nil)
				userRepo.On("CreateUser", mock.Anything, tx, mock.Anything).Return(nil)

				walletRepo.On("CreateWallet", mock.Anything, tx, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Empty Name",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "",
				email:    "user@example.com",
				password: "password123",
				role:     models.Customer,
			},
			mockSetup:     func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {},
			expectedError: errors.ErrNameEmpty,
		},
		{
			name: "Empty Email",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "User",
				email:    "",
				password: "password123",
				role:     models.Customer,
			},
			mockSetup:     func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {},
			expectedError: errors.ErrEmailEmpty,
		},
		{
			name: "Empty Password",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "User",
				email:    "user@example.com",
				password: "",
				role:     models.Customer,
			},
			mockSetup:     func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {},
			expectedError: errors.ErrPasswordEmpty,
		},
		{
			name: "Password Too Short",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "User",
				email:    "user@example.com",
				password: "123",
				role:     models.Customer,
			},
			mockSetup:     func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {},
			expectedError: errors.ErrPasswordTooShort,
		},
		{
			name: "Email Already Exists",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "User",
				email:    "existing@example.com",
				password: "password123",
				role:     models.Customer,
			},
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userRepo.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(&models.User{ID: 1, Email: "existing@example.com"}, nil)
			},
			expectedError: errors.ErrUserAlreadyExists,
		},
		{
			name: "Username Taken",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "JohnDoe",
				email:    "johndoe@example.com",
				password: "securepass",
				role:     models.Customer,
			},
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				userRepo.On("GetUserByEmail", mock.Anything, "johndoe@example.com").Return((*models.User)(nil), nil)
				userRepo.On("GetUserByName", mock.Anything, "JohnDoe").Return(&models.User{ID: 2, Name: "JohnDoe"}, nil)
			},
			expectedError: errors.ErrUsernameTaken,
		},
		{
			name: "Create User Failed",
			input: struct {
				name     string
				email    string
				password string
				role     models.UserRole
			}{
				name:     "User",
				email:    "user@example.com",
				password: "password123",
				role:     models.Customer,
			},
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("BeginTransaction", mock.Anything).Return(tx)
				userRepo.On("GetUserByEmail", mock.Anything, "user@example.com").Return((*models.User)(nil), nil)
				userRepo.On("GetUserByName", mock.Anything, "User").Return((*models.User)(nil), nil)
				userRepo.On("CreateUser", mock.Anything, tx, mock.Anything).Return(errors.ErrCreatedUser)
			},
			expectedError: errors.ErrCreatedUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			mockWalletRepo := new(mocks.WalletRepository)
			mockTransaction := new(mocks.Transaction)

			userUsecase := usecase.NewUserUsecase(mockUserRepo, mockWalletRepo)

			tt.mockSetup(mockUserRepo, mockWalletRepo, mockTransaction)

			_, err := userUsecase.RegisterUser(ctx, tt.input.name, tt.input.email, tt.input.password, tt.input.role)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		mockSetup     func(*mocks.UserRepository, *mocks.WalletRepository)
		expectedError error
		expectedToken string
	}{
		{
			name:     "Success",
			email:    "johndoe@example.com",
			password: "securepass",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepass"), bcrypt.DefaultCost)
				user := &models.User{ID: 1, Email: "johndoe@example.com", Password: string(hashedPassword)}

				userRepo.On("GetUserByEmail", mock.Anything, "johndoe@example.com").Return(user, nil)
			},
			expectedError: nil,
			expectedToken: "valid_token",
		},
		{
			name:     "User Not Found",
			email:    "notfound@example.com",
			password: "securepass",
			mockSetup: func(userRepo *mocks.UserRepository, wallet *mocks.WalletRepository) {
				userRepo.On("GetUserByEmail", mock.Anything, "notfound@example.com").Return((*models.User)(nil), errors.ErrUserNotFound)
			},
			expectedError: errors.ErrUserNotFound,
			expectedToken: "",
		},
		{
			name:     "Invalid Password",
			email:    "johndoe@example.com",
			password: "wrongpass",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepass"), bcrypt.DefaultCost)
				user := &models.User{ID: 1, Email: "johndoe@example.com", Password: string(hashedPassword)}

				userRepo.On("GetUserByEmail", mock.Anything, "johndoe@example.com").Return(user, nil)
			},
			expectedError: errors.ErrInvalidCredentials,
			expectedToken: "",
		},
		{
			name:     "Token Generation Failed",
			email:    "johndoe@example.com",
			password: "securepass",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("securepass"), bcrypt.DefaultCost)
				user := &models.User{ID: 1, Email: "johndoe@example.com", Password: string(hashedPassword)}

				userRepo.On("GetUserByEmail", mock.Anything, "johndoe@example.com").Return(user, nil)
			},
			expectedError: errors.ErrTokenGeneration,
			expectedToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockUserRepo := new(mocks.UserRepository)
			walletRepo := new(mocks.WalletRepository)

			userUsecase := usecase.NewUserUsecase(mockUserRepo, walletRepo)

			tt.mockSetup(mockUserRepo, walletRepo)

			patch := monkey.Patch(JWT.GenerateJWT, func(user *models.User) (string, error) {
				if tt.expectedError == errors.ErrTokenGeneration {
					return "", errors.ErrTokenGeneration
				}
				return "valid_token", nil
			})
			defer patch.Unpatch()

			token, err := userUsecase.LoginUser(ctx, tt.email, tt.password)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestGetInfoUser(t *testing.T) {
	tests := []struct {
		name                 string
		mockSetup            func(*mocks.UserRepository, *mocks.WalletRepository)
		expectedError        error
		expectedUser         *models.User
		expectedBalance      float64
		expectedTransactions *[]models.Transaction
	}{
		{
			name: "Success",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				userID := int64(1)
				user := &models.User{ID: userID, Name: "Alice", Email: "alice@example.com"}
				balance := 500.0
				transactions := []models.Transaction{
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
			expectedTransactions: &[]models.Transaction{
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
			name: "Wallet Not Found",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				userID := int64(1)
				user := &models.User{ID: userID, Name: "Alice", Email: "alice@example.com"}

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(0.0, errors.ErrCheckBalance)
			},
			expectedError: errors.ErrCheckBalance,
		},
		{
			name: "Transaction Fetch Failed",
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository) {
				userID := int64(1)
				user := &models.User{ID: userID, Name: "Alice", Email: "alice@example.com"}
				balance := 500.0

				userRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
				walletRepo.On("GetBalance", mock.Anything, userID).Return(balance, nil)
				walletRepo.On("GetTransactions", mock.Anything, userID).Return([]models.Transaction{}, errors.ErrCheckTransaction)
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
