package usecase_test

import (
	"bou.ke/monkey"
	"computer-club/internal/models"
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

func TestRegister(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
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
				email    string
				password string
				role     models.UserRole
			}{
				email:    "johndoe@example.com",
				password: "securepass",
				role:     models.Customer,
			},
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("BeginTransaction", mock.Anything).Return(tx)
				userRepo.On("GetUserByEmail", mock.Anything, "johndoe@example.com").Return((*models.User)(nil), nil)
				userRepo.On("CreateUser", mock.Anything, tx, mock.Anything).Return(nil)
				walletRepo.On("CreateWallet", mock.Anything, tx, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Empty Email",
			input: struct {
				email    string
				password string
				role     models.UserRole
			}{
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
				email    string
				password string
				role     models.UserRole
			}{
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
				email    string
				password string
				role     models.UserRole
			}{
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
				email    string
				password string
				role     models.UserRole
			}{
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
			name: "Create User Failed",
			input: struct {
				email    string
				password string
				role     models.UserRole
			}{
				email:    "user@example.com",
				password: "password123",
				role:     models.Customer,
			},
			mockSetup: func(userRepo *mocks.UserRepository, walletRepo *mocks.WalletRepository, tx *mocks.Transaction) {
				tx.On("Rollback").Return(nil)
				tx.On("Commit").Return(nil)

				userRepo.On("BeginTransaction", mock.Anything).Return(tx)
				userRepo.On("GetUserByEmail", mock.Anything, "user@example.com").Return((*models.User)(nil), nil)
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

			authUsecase := usecase.NewAuthUsecase(mockUserRepo, mockWalletRepo)

			tt.mockSetup(mockUserRepo, mockWalletRepo, mockTransaction)

			_, err := authUsecase.Register(ctx, tt.input.email, tt.input.password, tt.input.role)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestLogin(t *testing.T) {
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

			authUsecase := usecase.NewAuthUsecase(mockUserRepo, walletRepo)

			tt.mockSetup(mockUserRepo, walletRepo)

			patch := monkey.Patch(JWT.GenerateJWT, func(user *models.User) (string, error) {
				if tt.expectedError == errors.ErrTokenGeneration {
					return "", errors.ErrTokenGeneration
				}
				return "valid_token", nil
			})
			defer patch.Unpatch()

			token, err := authUsecase.Login(ctx, tt.email, tt.password)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}
