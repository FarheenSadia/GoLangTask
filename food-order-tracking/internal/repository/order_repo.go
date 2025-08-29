package repository

import (
	"food-order-tracking/internal/models"
	"time"

	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepo) AddEvent(orderID, event string) error {
	ev := models.OrderEvent{
		OrderID:   orderID,
		Event:     event,
		Timestamp: time.Now(),
	}
	return r.db.Create(&ev).Error
}

func (r *OrderRepo) UpdateStatus(orderID, status string) error {
	return r.db.Model(&models.Order{}).Where("order_id = ?", orderID).
		Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

func (r *OrderRepo) GetOrder(orderID string) (models.Order, []models.OrderEvent, error) {
	var order models.Order
	var events []models.OrderEvent
	err := r.db.First(&order, "order_id = ?", orderID).Error
	if err != nil {
		return order, nil, err
	}
	r.db.Where("order_id = ?", orderID).Order("timestamp asc").Find(&events)
	return order, events, nil
}
