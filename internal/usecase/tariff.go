package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"computer-club/pkg/errors"
	"context"
)

type TariffService interface {
	GetTariff(ctx context.Context) ([]models.Tariff, error)
	GetTariffByID(ctx context.Context, id int64) (*models.Tariff, error)
	DeleteTariffByID(ctx context.Context, id int64) error
	CreateTariff(ctx context.Context, id int64, name string, price float64, duration int64) (*models.Tariff, error)
}

type TariffUsecase struct {
	tariffRepository repository.TariffRepository
}

func (u *TariffUsecase) CreateTariff(ctx context.Context, id int64, name string, price float64, duration int64) (*models.Tariff, error) {
	if name == "" || price <= 0 || duration <= 0 || id <= 0 {
		return nil, errors.ErrInvalidInputDataTariff
	}

	_, err := u.tariffRepository.GetTariffByID(ctx, id)
	if err == nil {
		return nil, errors.ErrTariffWithIdAlreadyExist
	}

	tariff := &models.Tariff{
		ID:       id,
		Name:     name,
		Price:    price,
		Duration: duration,
	}

	return u.tariffRepository.CreateTariff(ctx, tariff)

}

func NewTariffUsecase(tariffRepository repository.TariffRepository) TariffService {
	return &TariffUsecase{tariffRepository: tariffRepository}
}

func (u *TariffUsecase) GetTariff(ctx context.Context) ([]models.Tariff, error) {
	return u.tariffRepository.GetTariff(ctx)
}

func (u *TariffUsecase) GetTariffByID(ctx context.Context, id int64) (*models.Tariff, error) {
	return u.tariffRepository.GetTariffByID(ctx, id)
}

func (u *TariffUsecase) DeleteTariffByID(ctx context.Context, id int64) error {
	tariff, err := u.tariffRepository.GetTariffByID(ctx, id)
	if err != nil {
		return err
	}
	return u.tariffRepository.DeleteTariffByID(ctx, tariff)
}
