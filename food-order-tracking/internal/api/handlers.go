package api

import (
	"food-order-tracking/internal/models"
	"food-order-tracking/internal/repository"

	//"food-order-tracking/internal/kafka"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	repo *repository.OrderRepo
}

func NewHandler(repo *repository.OrderRepo) *Handler {
	return &Handler{repo: repo}
}

// POST /orders
func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	var req models.Order
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	// assign metadata
	req.OrderID = "ORD-" + uuid.New().String()
	req.Status = "PLACED"
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	// persist order
	if err := h.repo.CreateOrder(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot save order"})
	}

	// produce Kafka event
	// err := kafka.ProduceOrder(kafka.OrderMessage{
	// 	OrderID:      req.OrderID,
	// 	CustomerName: req.CustomerName,
	// 	Address:      req.Address,
	// 	Item:         req.Item,
	// 	Size:         req.Size,
	// 	Status:       req.Status,
	// 	CreatedAt:    req.CreatedAt,
	// })
	// if err != nil {
	// 	// just log the error, don’t block response
	// 	// (you can later use a logger package)
	// 	println("failed to publish Kafka event:", err.Error())
	// }

	// save event in DB
	if err := h.repo.AddEvent(req.OrderID, "CREATED"); err != nil {
		println("failed to log event:", err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"order_id": req.OrderID})
}

// GET /orders/:order_id
func (h *Handler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("order_id")

	order, events, err := h.repo.GetOrder(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}

	return c.JSON(fiber.Map{
		"order":  order,
		"events": events,
	})
}

// GET /health
func (h *Handler) Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
