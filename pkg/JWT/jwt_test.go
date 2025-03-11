package JWT_test

import (
	"computer-club/internal/repository/models"
	"computer-club/pkg/JWT"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name          string
		user          *models.User
		mockSetup     func()
		expectedError error
	}{
		{
			name: "Success",
			user: &models.User{
				ID:   1,
				Role: "customer",
			},
			mockSetup: func() {
				os.Setenv("JWT_SECRET", "test_secret")
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			token, err := JWT.GenerateJWT(tt.user)

			if tt.expectedError == nil {
				assert.NotEmpty(t, token)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, token)
				assert.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}
