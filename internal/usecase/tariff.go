package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"context"
)

type TariffService interface {
	GetTariff(ctx context.Context) ([]models.Tariff, error)
	GetTariffByID(ctx context.Context, id int64) (*models.Tariff, error)
}
type TariffUsecase struct {
	tariffRepository repository.TariffRepository
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
