package models

import "time"

type SessionStatus string

const (
	Active   SessionStatus = "active"
	Finished SessionStatus = "finished"
)

// Session - модель сессии компьютера
type Session struct {
	ID        int64         `json:"id" gorm:"primaryKey"`
	UserID    int64         `json:"user_id" gorm:"index"`
	PCNumber  int           `json:"pc_number"`
	TariffID  int64         `json:"tariff_id"`
	Status    SessionStatus `json:"status"`
	StartTime time.Time     `json:"start_time"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
}
