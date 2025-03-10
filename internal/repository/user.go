package repository

import (
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByName(ctx context.Context, name string) (*models.User, error)
}

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errors.ErrCreatedUser
	}

	return nil
}

func (r *PostgresUserRepo) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	var user models.User
	err := r.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}
