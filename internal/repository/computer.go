package repository

import (
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type ComputerRepository interface {
	GetComputers(ctx context.Context) ([]models.Computer, error)
	UpdateStatus(ctx context.Context, number int, free models.ComputerStatus) error
	IsComputerAvailable(ctx context.Context, number int) (bool, error)
}
type PostgresComputerRepo struct {
	db *gorm.DB
}

func NewComputerRepository(db *gorm.DB) ComputerRepository {
	return &PostgresComputerRepo{db: db}
}

func (r *PostgresComputerRepo) UpdateStatus(ctx context.Context, number int, free models.ComputerStatus) error {
	var computer models.Computer
	if err := r.db.WithContext(ctx).First(&computer, "id = ?", number).Error; err != nil {
		return errors.ErrFindComputer
	}
	if err := r.db.WithContext(ctx).Model(&computer).Update("status", free).Error; err != nil {
		return errors.ErrUpdateComputer
	}
	return nil
}
func (r *PostgresComputerRepo) GetComputers(ctx context.Context) ([]models.Computer, error) {
	var computers []models.Computer
	if err := r.db.WithContext(ctx).Find(&computers).Error; err != nil {
		return nil, errors.ErrFindComputer
	}
	return computers, nil
}
func (r *PostgresComputerRepo) IsComputerAvailable(ctx context.Context, number int) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Computer{}).
		Where("pc_number = ? AND status = ?", number, models.Free).
		Count(&count).
		Error; err != nil {
		return false, errors.ErrPCBusy
	}
	return count > 0, nil
}
