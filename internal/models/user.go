package models

// UserRole - роли пользователей
type UserRole string

const (
	Admin    UserRole = "admin"
	Customer UserRole = "customer"
)

// User - модель пользователя
type User struct {
	ID           int64         `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name"`
	Role         string        `json:"role"`
	Email        string        `json:"email"`
	Password     string        `json:"-"`
	Wallet       *Wallet       `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
	Sessions     []Session     `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
	Transactions []Transaction `json:"-" gorm:"constraint:OnDelete:CASCADE;foreignKey:UserID"`
}
