package repository

import (
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type ComputerRepository interface {
	GetComputers(ctx context.Context) ([]models.Computer, error)
	IsComputerAvailable(ctx context.Context, tx Transaction, number int) (bool, error)
	MarkComputerFree(ctx context.Context, tx Transaction, number int) error
}
type PostgresComputerRepo struct {
	db *gorm.DB
}

func NewComputerRepository(db *gorm.DB) ComputerRepository {
	return &PostgresComputerRepo{db: db}
}

func (r *PostgresComputerRepo) GetComputers(ctx context.Context) ([]models.Computer, error) {
	var computers []models.Computer
	if err := r.db.WithContext(ctx).Find(&computers).Error; err != nil {
		return nil, errors.ErrFindComputer
	}
	return computers, nil
}

func (r *PostgresComputerRepo) IsComputerAvailable(ctx context.Context, tx Transaction, number int) (bool, error) {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}

	var computer models.Computer
	if err := db.WithContext(ctx).
		Where("pc_number = ?", number).First(&computer).Error; err != nil {
		return false, errors.ErrComputerNotFound
	}
	if computer.Status == models.Busy {
		return true, nil
	}
	return false, nil
}

func (r *PostgresComputerRepo) MarkComputerFree(ctx context.Context, tx Transaction, pcNumber int) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	err := db.WithContext(ctx).Model(&models.Computer{}).
		Where("pc_number = ?", pcNumber).
		Update("status", models.Free).Error
	if err != nil {
		return errors.ErrUpdateComputerStatus
	}
	return nil
}
