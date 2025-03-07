package repository

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"gorm.io/gorm"
)

type ComputerRepository interface {
	GetComputers() ([]models.Computer, error)
}
type PostgresComputerRepo struct {
	db *gorm.DB
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
