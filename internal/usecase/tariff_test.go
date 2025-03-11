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

func TestGetTariff(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*mocks.TariffRepository)
		expectedError error
		expectedData  []models.Tariff
	}{
		{
			name: "Success",
			mockSetup: func(tariffRepo *mocks.TariffRepository) {
				tariffs := []models.Tariff{
					{ID: 1, Name: "Basic", Price: 100},
					{ID: 2, Name: "Premium", Price: 200},
				}
				tariffRepo.On("GetTariff", mock.Anything).Return(tariffs, nil)
			},
			expectedError: nil,
			expectedData: []models.Tariff{
				{ID: 1, Name: "Basic", Price: 100},
				{ID: 2, Name: "Premium", Price: 200},
			},
		},
		{
			name: "Repository Error",
			mockSetup: func(tariffRepo *mocks.TariffRepository) {
				tariffRepo.On("GetTariff", mock.Anything).Return(nil, errors.ErrFindTariffs)
			},
			expectedError: errors.ErrFindTariffs,
			expectedData:  nil,
		},
		{
			name: "No Tariffs Found",
			mockSetup: func(tariffRepo *mocks.TariffRepository) {
				tariffRepo.On("GetTariff", mock.Anything).Return([]models.Tariff{}, nil)
			},
			expectedError: nil,
			expectedData:  []models.Tariff{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockTariffRepo := new(mocks.TariffRepository)

			tariffUsecase := usecase.NewTariffUsecase(mockTariffRepo)

			tt.mockSetup(mockTariffRepo)

			result, err := tariffUsecase.GetTariff(ctx)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedData, result)
		})
	}
}

func TestGetTariffByID(t *testing.T) {
	tests := []struct {
		name          string
		tariffID      int64
		mockSetup     func(*mocks.TariffRepository, int64)
		expectedError error
		expectedData  *models.Tariff
	}{
		{
			name:     "Success",
			tariffID: 1,
			mockSetup: func(tariffRepo *mocks.TariffRepository, tariffID int64) {
				tariff := &models.Tariff{ID: tariffID, Name: "Basic", Price: 100}
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return(tariff, nil)
			},
			expectedError: nil,
			expectedData:  &models.Tariff{ID: 1, Name: "Basic", Price: 100},
		},
		{
			name:     "Tariff Not Found",
			tariffID: 99,
			mockSetup: func(tariffRepo *mocks.TariffRepository, tariffID int64) {
				tariffRepo.On("GetTariffByID", mock.Anything, tariffID).Return((*models.Tariff)(nil), errors.ErrFindTariffByID)
			},
			expectedError: errors.ErrFindTariffByID,
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockTariffRepo := new(mocks.TariffRepository)

			tariffUsecase := usecase.NewTariffUsecase(mockTariffRepo)

			tt.mockSetup(mockTariffRepo, tt.tariffID)

			result, err := tariffUsecase.GetTariffByID(ctx, tt.tariffID)

			assert.ErrorIs(t, err, tt.expectedError)
			assert.Equal(t, tt.expectedData, result)
		})
	}
}
