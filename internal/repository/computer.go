package repository

import (
	"computer-club/internal/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type ComputerRepository interface {
	GetComputers(ctx context.Context) ([]models.Computer, error)
	UpdateStatus(ctx context.Context, number int, free models.ComputerStatus) error
}
type PostgresComputerRepo struct {
	db *gorm.DB
}

func NewComputerRepository(db *gorm.DB) ComputerRepository {
	return &PostgresComputerRepo{db: db}
}

func (r *PostgresComputerRepo) UpdateStatus(ctx context.Context, number int, free models.ComputerStatus) error {
	var computer models.Computer
	if err := r.db.First(&computer, "id = ?", number).Error; err != nil {
		return errors.ErrFindComputer
	}
	if err := r.db.Model(&computer).Update("status", free).Error; err != nil {
		return errors.ErrUpdateComputer
	}
	return nil
}
func (r *PostgresComputerRepo) GetComputers(ctx context.Context) ([]models.Computer, error) {
	var computers []models.Computer
	if err := r.db.Find(&computers).Error; err != nil {
		return nil, errors.ErrFindComputer
	}
	return computers, nil
}
