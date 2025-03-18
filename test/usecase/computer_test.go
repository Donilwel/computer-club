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

func TestGetComputersStatus(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.ComputerRepository)
		expectedError error
		expectedData  []*models.Computer
	}{
		{
			name: "Success",
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computers := []*models.Computer{
					{ID: 1, PCNumber: 101, Status: models.Free},
					{ID: 2, PCNumber: 102, Status: models.Busy},
				}
				computerRepo.On("GetComputers", mock.Anything).Return(computers, nil)
			},
			expectedError: nil,
			expectedData: []*models.Computer{
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

func TestDeleteComputer(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		mockSetup     func(*mocks.ComputerRepository)
		expectedError error
	}{
		{
			name: "Success",
			id:   1,
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computerRepo.On("GetComputerByID", mock.Anything, 1).Return(&models.Computer{ID: 1, Status: models.Free}, nil)
				computerRepo.On("DeleteComputer", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Computer not found",
			id:   2,
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computerRepo.On("GetComputerByID", mock.Anything, 2).Return(nil, errors.ErrComputerNotFound)
			},
			expectedError: errors.ErrComputerNotFound,
		},
		{
			name: "Computer is busy",
			id:   3,
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computerRepo.On("GetComputerByID", mock.Anything, 3).Return(&models.Computer{ID: 3, Status: models.Busy}, nil)
			},
			expectedError: errors.ErrPCBusy,
		},
		{
			name: "Error Delete Computer",
			id:   4,
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computerRepo.On("GetComputerByID", mock.Anything, 4).Return(&models.Computer{ID: 4, Status: models.Free}, nil)
				computerRepo.On("DeleteComputer", mock.Anything, mock.Anything).Return(errors.ErrDeleteComputer)
			},
			expectedError: errors.ErrDeleteComputer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockComputerRepo := new(mocks.ComputerRepository)

			computerUsecase := usecase.NewComputerUsecase(mockComputerRepo)

			tt.mockSetup(mockComputerRepo)

			err := computerUsecase.DeleteComputer(ctx, tt.id)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestAddComputer(t *testing.T) {
	tests := []struct {
		name             string
		mockSetup        func(*mocks.ComputerRepository)
		expectedError    error
		expectedComputer *models.Computer
	}{
		{
			name: "Success",
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computerRepo.On("AddComputer", mock.Anything).Return(&models.Computer{
					ID:     1,
					Status: models.Free,
				}, nil)
			},
			expectedError: nil,
			expectedComputer: &models.Computer{
				ID:     1,
				Status: models.Free,
			},
		},
		{
			name: "Failed to add computer",
			mockSetup: func(computerRepo *mocks.ComputerRepository) {
				computerRepo.On("AddComputer", mock.Anything).Return(nil, errors.ErrCreateComputer)
			},
			expectedError:    errors.ErrCreateComputer,
			expectedComputer: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockComputerRepo := new(mocks.ComputerRepository)

			computerUsecase := usecase.NewComputerUsecase(mockComputerRepo)

			tt.mockSetup(mockComputerRepo)

			computer, err := computerUsecase.AddComputer(ctx)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedComputer, computer)
		})
	}
}
