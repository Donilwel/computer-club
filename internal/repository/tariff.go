package repository

import (
	"computer-club/internal/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type TariffRepository interface {
	GetTariff(ctx context.Context) ([]models.Tariff, error)
	GetTariffByID(ctx context.Context, id int64) (*models.Tariff, error)
	DeleteTariffByID(ctx context.Context, tariff *models.Tariff) error
	CreateTariff(context.Context, *models.Tariff) (*models.Tariff, error)
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

func (r *TariffRepositoryPostgres) CreateTariff(ctx context.Context, tariff *models.Tariff) (*models.Tariff, error) {
	err := r.db.WithContext(ctx).Create(tariff).Error
	if err != nil {
		return nil, errors.ErrCreateTariff
	}
	return tariff, nil
}

func (r *TariffRepositoryPostgres) DeleteTariffByID(ctx context.Context, tariff *models.Tariff) error {
	if err := r.db.WithContext(ctx).Delete(tariff).Error; err != nil {
		return errors.ErrDeleteTariff
	}
	return nil
}
