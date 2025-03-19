package models

type Wallet struct {
	ID      int64   `json:"id" gorm:"primaryKey"`
	UserID  int64   `json:"user_id" gorm:"index"`
	Balance float64 `json:"balance"`
}
