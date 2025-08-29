package models

import "time"

type Order struct {
	OrderID      string    `gorm:"primaryKey" json:"order_id"`
	CustomerName string    `json:"customer_name"`
	Address      string    `json:"address"`
	Item         string    `json:"item"`
	Size         string    `json:"size"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type OrderEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   string    `json:"order_id"`
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Meta      string    `json:"meta"`
}
