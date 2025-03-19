package models

import "time"

type SessionStatus string

const (
	Active   SessionStatus = "active"
	Finished SessionStatus = "finished"
)

// Session - модель сессии компьютера
type Session struct {
	ID         int64      `json:"id" gorm:"primaryKey"`
	UserID     int64      `json:"user_id" gorm:"index"`                                        // Внешний ключ
	User       User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`     // Пользователь
	ComputerID int64      `json:"computer_id" gorm:"index"`                                    // Внешний ключ
	Computer   Computer   `gorm:"foreignKey:ComputerID;constraint:OnDelete:CASCADE;" json:"-"` // Компьютер
	TariffID   int64      `json:"tariff_id" gorm:"index"`
	Tariff     Tariff     `gorm:"foreignKey:TariffID;constraint:OnDelete:SET NULL;" json:"-"` // Связь с тарифом
	StartTime  time.Time  `json:"start_time"`
	EndTime    *time.Time `json:"end_time"` // Может быть NULL, если сессия активна
}
