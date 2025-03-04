package domain

import "context"

type User struct {
	ID       string
	Email    string
	Password string
}

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
