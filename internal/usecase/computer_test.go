package usecase_test

import (
	"computer-club/internal/repository/models"
	"computer-club/internal/usecase"
	"computer-club/mocks"
	"computer-club/pkg/errors"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGetComputersStatus(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.ComputerRepository)
		expectedError error
		expectedData  []models.Computer
	}{
		{
			name: "Success",
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computers := []models.Computer{
					{ID: 1, PCNumber: 101, Status: models.Free},
					{ID: 2, PCNumber: 102, Status: models.Busy},
				}
				computerRepo.On("GetComputers", mock.Anything).Return(computers, nil)
			},
			expectedError: nil,
			expectedData: []models.Computer{
				{ID: 1, PCNumber: 101, Status: models.Free},
				{ID: 2, PCNumber: 102, Status: models.Busy},
			},
		},
		{
			name: "Error Find Computer",
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computerRepo.On("GetComputers", mock.Anything).Return(nil, errors.ErrFindComputer)
			},
			expectedError: errors.ErrFindComputer,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockComputerRepo := new(mocks.ComputerRepository)

			computerUsecase := usecase.NewComputerUsecase(mockComputerRepo)

			tt.mockSetup(mockComputerRepo)

			result, err := computerUsecase.GetComputersStatus(ctx)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedData, result)
		})
	}
}
