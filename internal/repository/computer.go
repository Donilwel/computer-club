package repository

import (
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
	err := r.db.Find(&computers).Error
	return computers, err
}
