package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"payments-service/database"
	"payments-service/models"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx  = context.Background()
	rcli *redis.Client
)

type OrderEvent struct {
	OrderID int     `json:"order_id"`
	UserID  int     `json:"user_id"`
	Amount  float64 `json:"amount"`
}

func Client() *redis.Client {
	if rcli != nil {
		return rcli
	}
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}
	rcli = redis.NewClient(&redis.Options{Addr: addr})
	return rcli
}

func StartConsumer(rdb *redis.Client) {
	sub := rdb.Subscribe(ctx, "orders")
	defer sub.Close()

	ch := sub.Channel()
	log.Println("Payments consumer subscribed to 'orders'")

	for msg := range ch {
		var event OrderEvent
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Println("parse event error:", err)
			continue
		}

		rand.Seed(time.Now().UnixNano())
		outcome := rand.Intn(100)
		status := "failed"
		if outcome%2 == 0 {
			status = "success"
		}

		payment := models.Payment{
			OrderID: event.OrderID,
			Amount:  event.Amount,
			Status:  status,
		}
		if err := database.DB.Create(&payment).Error; err != nil {
			log.Println("payment insert error:", err)
			continue
		}

		base := os.Getenv("ORDERS_SERVICE_BASE")
		if base == "" {
			base = "http://users-orders-service:8080"
		}
		url := fmt.Sprintf("%s/orders/%d/status/%s", base, event.OrderID, status)

		resp, err := http.Post(url, "application/json", http.NoBody)
		if err != nil {
			log.Printf("order status update failed (order %d): %v", event.OrderID, err)
			continue
		}
		_ = resp.Body.Close()

		log.Printf("Processed payment for Order %d -> %s (payment_id=%d)", event.OrderID, status, payment.ID)
	}
}
