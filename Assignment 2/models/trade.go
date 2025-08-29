package models

import (
	"time"

	"gorm.io/gorm"
)

type Trade struct {
	ID        uint    `gorm:"primaryKey"`
	Symbol    string  `gorm:"not null"`
	Side      string  `gorm:"not null"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
