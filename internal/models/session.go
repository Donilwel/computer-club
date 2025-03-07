package models

import "time"

// Session - модель сессии компьютера
type Session struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	PCNumber  int        `json:"pc_number"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}
