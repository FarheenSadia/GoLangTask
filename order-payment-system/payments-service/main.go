package main

import (
	"log"
	"os"
	"payments-service/consumer"
	"payments-service/database"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}

	// web server (health)
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("ok") })

	// start consumer
	rdb := consumer.Client()
	go consumer.StartConsumer(rdb)

	port := getenv("APP_PORT", "8081")
	go func() {
		log.Printf("Payments Service HTTP on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatal(err)
		}
	}()

	// block main
	select {}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
