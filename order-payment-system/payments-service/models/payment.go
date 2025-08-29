package models

type Payment struct {
	ID      int     `gorm:"primaryKey" json:"id"`
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}
