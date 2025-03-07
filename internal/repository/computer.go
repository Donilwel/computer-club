package repository

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"gorm.io/gorm"
)

type ComputerRepository interface {
	GetComputers() ([]models.Computer, error)
	UpdateStatus(number int, free models.ComputerStatus) error
}
type PostgresComputerRepo struct {
	db *gorm.DB
}

func (r *PostgresComputerRepo) UpdateStatus(number int, free models.ComputerStatus) error {
	var computer models.Computer
	if err := r.db.First(&computer, "id = ?", number).Error; err != nil {
		return errors.ErrFindComputer
	}
	if err := r.db.Model(&computer).Update("status", free).Error; err != nil {
		return errors.ErrUpdateComputer
	}
	return nil
}

func NewComputerRepository(db *gorm.DB) ComputerRepository {
	return &PostgresComputerRepo{db: db}
}

func (r *PostgresComputerRepo) GetComputers() ([]models.Computer, error) {
	var computers []models.Computer
	if err := r.db.Find(&computers).Error; err != nil {
		return nil, errors.ErrFindComputer
	}
	return computers, nil
}
