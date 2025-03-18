package repository

import (
	"computer-club/internal/models"
	"computer-club/pkg/errors"
	"computer-club/pkg/generate"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx Transaction, user *models.User) error
	GetUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByName(ctx context.Context, name string) (*models.User, error)
	BeginTransaction(ctx context.Context) Transaction
}

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) UserRepository {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) BeginTransaction(ctx context.Context) Transaction {
	return &GormTransaction{tx: r.db.WithContext(ctx).Begin()}
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, tx Transaction, user *models.User) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}

	user.Name = generate.GenerateUsername(db)

	if err := db.WithContext(ctx).Create(user).Error; err != nil {
		return errors.ErrCreatedUser
	}

	return nil
}

func (r *PostgresUserRepo) GetUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User

	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, errors.ErrFindUsers
	}

	if len(users) == 0 {
		return nil, errors.ErrZeroUsers
	}
	return users, nil
}

func (r *PostgresUserRepo) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)
	if result.Error != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetUserByName(ctx context.Context, name string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, errors.ErrFindUser
	}
	return &user, nil
}
