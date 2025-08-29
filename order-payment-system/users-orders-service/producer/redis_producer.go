package producer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	once   sync.Once
	client *redis.Client
	ctx    = context.Background()
)

func Client() *redis.Client {
	once.Do(func() {
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "redis:6379"
		}
		client = redis.NewClient(&redis.Options{Addr: addr})
	})
	return client
}

type OrderEvent struct {
	OrderID int     `json:"order_id"`
	UserID  int     `json:"user_id"`
	Amount  float64 `json:"amount"`
}

func PublishOrder(rdb *redis.Client, event OrderEvent) {
	data, _ := json.Marshal(event)
	if err := rdb.Publish(ctx, "orders", data).Err(); err != nil {
		log.Println("redis publish error:", err)
	}
}
