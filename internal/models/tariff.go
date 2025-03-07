package models

import "time"

type Tariff struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Duration  int64     `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}
