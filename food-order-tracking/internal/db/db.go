package db

import (
	"log"
	"os"

	"food-order-tracking/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dsn := os.Getenv("PG_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect DB: ", err)
	}
	db.AutoMigrate(&models.Order{}, &models.OrderEvent{})
	return db
}
