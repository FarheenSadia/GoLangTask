package main

import (
	"log"
	"os"
	"users-orders-service/database"
	"users-orders-service/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	routes.SetupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Users & Orders Service running on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
