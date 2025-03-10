package repository

import (
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type TariffRepository interface {
	GetTariff(ctx context.Context) ([]models.Tariff, error)
	GetTariffByID(ctx context.Context, id int64) (*models.Tariff, error)
}

type TariffRepositoryPostgres struct {
	db *gorm.DB
}

func NewTariffRepositoryPostgres(db *gorm.DB) TariffRepository {
	return &TariffRepositoryPostgres{db: db}
}

func (r *TariffRepositoryPostgres) GetTariff(ctx context.Context) ([]models.Tariff, error) {
	var tariffs []models.Tariff
	err := r.db.WithContext(ctx).Find(&tariffs).Error
	if err != nil {
		return nil, errors.ErrFindTariffs
	}
	return tariffs, nil
}

func (r *TariffRepositoryPostgres) GetTariffByID(ctx context.Context, id int64) (*models.Tariff, error) {
	var tariff models.Tariff
	err := r.db.WithContext(ctx).First(&tariff, id).Error
	if err != nil {
		return nil, errors.ErrFindTariffByID
	}
	return &tariff, nil
}
