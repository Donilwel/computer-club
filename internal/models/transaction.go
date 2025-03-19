package models

import "time"

type TransactionType string

const (
	Add TransactionType = "add"
	Buy TransactionType = "buy"
)

type Transaction struct {
	ID        int64           `json:"id" gorm:"primaryKey"`
	UserID    int64           `json:"user_id" gorm:"index"`                                    // Внешний ключ
	User      User            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"` // Связь с пользователем
	Amount    float64         `json:"amount"`
	TariffID  int64           `json:"tariff_id"`                                                  // Внешний ключ
	Tariff    Tariff          `gorm:"foreignKey:TariffID;constraint:OnDelete:SET NULL;" json:"-"` // Связь с тарифом
	Type      TransactionType `json:"type"`
	CreatedAt time.Time       `json:"created_at" gorm:"index"`
}
