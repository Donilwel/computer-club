package models

import "time"

type TransactionType string

const (
	Add TransactionType = "add"
	Buy TransactionType = "buy"
)

type Transaction struct {
	ID        int64           `json:"id" gorm:"primaryKey"`
	UserID    int64           `json:"user_id" gorm:"index"`
	Amount    float64         `json:"amount"`
	TariffID  int64           `json:"tariff_id"`
	Type      TransactionType `json:"type"`
	CreatedAt time.Time       `json:"created_at" gorm:"index"`
}
