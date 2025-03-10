package models

type Wallet struct {
	ID      int64   `json:"id"`
	UserID  int64   `json:"user_id"`
	Balance float64 `json:"balance"`
}
