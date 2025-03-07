package models

import "time"

// Session - модель сессии компьютера
type Session struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	PCNumber  int        `json:"pc_number"`
	TariffID  int64      `json:"tariff_id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}
