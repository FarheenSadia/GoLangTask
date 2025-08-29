package handlers

import (
	"strconv"
	"users-orders-service/database"
	"users-orders-service/models"
	"users-orders-service/producer"

	"github.com/gofiber/fiber/v2"
)

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if order.UserID == 0 || order.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_id and amount required"})
	}

	order.Status = "pending"
	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	producer.PublishOrder(producer.Client(), producer.OrderEvent{
		OrderID: order.ID,
		UserID:  order.UserID,
		Amount:  order.Amount,
	})

	return c.Status(fiber.StatusCreated).JSON(order)
}

func UpdateOrderStatus(c *fiber.Ctx) error {
	idParam := c.Params("id")
	status := c.Params("status")
	if status != "success" && status != "failed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "status must be success or failed"})
	}
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order id"})
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order not found"})
	}
	order.Status = status
	if err := database.DB.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(order)
}
