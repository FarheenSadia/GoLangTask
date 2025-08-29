package routes

import (
	"users-orders-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("ok") })

	app.Post("/users", handlers.CreateUser)
	app.Post("/orders", handlers.CreateOrder)
	app.Post("/orders/:id/status/:status", handlers.UpdateOrderStatus)
}
