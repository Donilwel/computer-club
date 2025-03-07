package repository

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByName(name string) (*models.User, error)
}

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.ErrCreatedUser
	}
	return nil
}

func (r *PostgresUserRepo) GetUserByID(id int64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetUserByName(name string) (*models.User, error) {
	var user models.User
	err := r.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}
