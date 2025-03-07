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

//type memoryUserRepo struct {
//	mu     sync.Mutex
//	users  map[int64]*models.User
//	lastID int64
//}
//
//func NewMemoryUserRepo() UserRepository {
//	return &memoryUserRepo{
//		users: make(map[int64]*models.User),
//	}
//}
//
//func (r *memoryUserRepo) CreateUser(user *models.User) error {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//	r.lastID++
//	user.ID = r.lastID
//	r.users[user.ID] = user
//	return nil
//}
//
//// GetUserByID получает пользователя по ID
//func (r *memoryUserRepo) GetUserByID(id int64) (*models.User, error) {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	user, exists := r.users[id]
//	if !exists {
//		return nil, fmt.Errorf("user not found")
//	}
//	return user, nil
//}

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
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
