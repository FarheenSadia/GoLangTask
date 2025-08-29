package models

type Order struct {
	ID     int     `gorm:"primaryKey" json:"id"`
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}
