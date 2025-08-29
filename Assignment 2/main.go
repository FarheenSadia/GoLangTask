package main

import (
	"log"
	"portfolio-tracker/config"
	"portfolio-tracker/models"
	"portfolio-tracker/services"
)

func main() {

	config.ConnectDB()

	if err := config.DB.AutoMigrate(&models.Trade{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	_ = services.AddTrade("INFY", "BUY", 100, 1400)
	_ = services.AddTrade("INFY", "SELL", 40, 1410)

	_ = services.GetNetPosition()
}
