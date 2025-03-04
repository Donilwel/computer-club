package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"computer-club/internal/domain"

	_ "github.com/lib/pq"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	var userID string
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password).Scan(&userID)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v", err)
		return "", err
	}
	return userID, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := "SELECT id, email, password FROM users WHERE email = $1"
	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Printf("Ошибка поиска пользователя: %v", err)
		return nil, err
	}
	return user, nil
}
