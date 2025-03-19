package models

// UserRole - роли пользователей
type UserRole string

const (
	Admin    UserRole = "admin"
	Customer UserRole = "customer"
)

// User - модель пользователя
type User struct {
	ID       int64     `json:"id" gorm:"primaryKey"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	Email    string    `json:"email"`
	Password string    `json:"-"`                                                       // Скрываем в JSON
	Wallet   Wallet    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"` // Связь один-к-одному
	Sessions []Session `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"` // Связь сессий с пользователем
}
