package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
)

type TariffService interface {
	GetTariff() ([]models.Tariff, error)
	GetTariffByID(id int64) (*models.Tariff, error)
}
type TariffUsecase struct {
	tariffRepository repository.TariffRepository
}

func NewTariffUsecase(tariffRepository repository.TariffRepository) *TariffUsecase {
	return &TariffUsecase{tariffRepository: tariffRepository}
}

func (u *TariffUsecase) GetTariff() ([]models.Tariff, error) {
	return u.tariffRepository.GetTariff()
}

func (u *TariffUsecase) GetTariffByID(id int64) (*models.Tariff, error) {
	return u.tariffRepository.GetTariffByID(id)
}
