package repository

import (
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type ComputerRepository interface {
	GetComputers(ctx context.Context) ([]*models.Computer, error)
	ChangeComputerStatus(ctx context.Context, tx Transaction, computer *models.Computer, status string) error
	GetComputerByID(ctx context.Context, id int) (*models.Computer, error)
	DeleteComputer(ctx context.Context, computer *models.Computer) error
	AddComputer(ctx context.Context) (*models.Computer, error)
}
type PostgresComputerRepo struct {
	db *gorm.DB
}

func (r *PostgresComputerRepo) GetComputerByID(ctx context.Context, id int) (*models.Computer, error) {
	var computer models.Computer
	err := r.db.WithContext(ctx).First(&computer, "pc_number = ?", id).Error
	if err != nil {
		return nil, errors.ErrComputerNotFound
	}
	return &computer, nil
}

func NewComputerRepository(db *gorm.DB) ComputerRepository {
	return &PostgresComputerRepo{db: db}
}

func (r *PostgresComputerRepo) GetComputers(ctx context.Context) ([]*models.Computer, error) {
	var computers []*models.Computer
	if err := r.db.WithContext(ctx).Find(&computers).Error; err != nil {
		return nil, errors.ErrFindComputer
	}
	return computers, nil
}

func (r *PostgresComputerRepo) ChangeComputerStatus(ctx context.Context, tx Transaction, computer *models.Computer, status string) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}

	if status != string(models.Busy) && status != string(models.Free) {
		return errors.ErrWrongComputerStatus
	}

	err := db.WithContext(ctx).Model(computer).Updates(models.Computer{Status: models.ComputerStatus(status)}).Error
	if err != nil {
		return errors.ErrUpdateComputerStatus
	}
	return nil
}

func (r *PostgresComputerRepo) DeleteComputer(ctx context.Context, computer *models.Computer) error {
	if err := r.db.WithContext(ctx).Model(computer).Delete(&computer).Error; err != nil {
		return errors.ErrDeleteComputer
	}
	return nil
}

func (r *PostgresComputerRepo) AddComputer(ctx context.Context) (*models.Computer, error) {
	computer := &models.Computer{}
	if err := r.db.WithContext(ctx).Create(computer).Error; err != nil {
		return nil, errors.ErrCreateComputer
	}
	computer.PCNumber = int(computer.ID)
	if err := r.db.WithContext(ctx).Save(computer).Error; err != nil {
		return nil, errors.ErrCreateComputer
	}
	return computer, nil
}
