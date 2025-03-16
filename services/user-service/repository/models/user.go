package models

// UserRole - роли пользователей
type UserRole string

const (
	Admin    UserRole = "admin"
	Customer UserRole = "customer"
)

// User - модель пользователя
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
