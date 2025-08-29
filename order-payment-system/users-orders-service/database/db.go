package database

import (
	"fmt"
	"os"
	"users-orders-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func dsn() string {
	host := getenv("PG_HOST", "users-db")
	user := getenv("PG_USER", "user")
	pass := getenv("PG_PASSWORD", "pass")
	db := getenv("PG_DB", "users_orders")
	port := getenv("PG_PORT", "5432")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, db, port)
}

func ConnectDB() error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("users db connect: %w", err)
	}
	if err := DB.AutoMigrate(&models.User{}, &models.Order{}); err != nil {
		return fmt.Errorf("users db migrate: %w", err)
	}
	return nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
