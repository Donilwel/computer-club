package models

import "time"

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

// Session - модель сессии компьютера
type Session struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	PCNumber  int        `json:"pc_number"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}
