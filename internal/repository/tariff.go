package repository

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"gorm.io/gorm"
)

type TariffRepository interface {
	GetTariff() ([]models.Tariff, error)
	GetTariffByID(id int64) (*models.Tariff, error)
}

type TariffRepositoryPostgres struct {
	db *gorm.DB
}

func NewTariffRepositoryPostgres(db *gorm.DB) *TariffRepositoryPostgres {
	return &TariffRepositoryPostgres{db: db}
}

func (r *TariffRepositoryPostgres) GetTariff() ([]models.Tariff, error) {
	var tariffs []models.Tariff
	err := r.db.Find(&tariffs).Error
	if err != nil {
		return nil, errors.ErrFindTariffs
	}
	return tariffs, nil
}

func (r *TariffRepositoryPostgres) GetTariffByID(id int64) (*models.Tariff, error) {
	var tariff models.Tariff
	err := r.db.First(&tariff, id).Error
	if err != nil {
		return nil, errors.ErrFindTariffByID
	}
	return &tariff, nil
}
