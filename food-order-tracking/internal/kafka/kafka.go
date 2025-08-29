package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	topic   = "orders.v1"
	brokers = []string{}
)

// Kafka message struct
type OrderMessage struct {
	OrderID      string    `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	Address      string    `json:"address"`
	Item         string    `json:"item"`
	Size         string    `json:"size"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func InitBrokers() {
	brokerStr := os.Getenv("KAFKA_BROKERS")
	if brokerStr == "" {
		brokerStr = "kafka:9092"
	}
	brokers = []string{brokerStr}
}

// Producer: send order message to Kafka
func ProduceOrder(msg OrderMessage) error {
	InitBrokers()
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	log.Printf("[Kafka] Producing order %s", msg.OrderID)

	return w.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(msg.OrderID),
			Value: data,
		},
	)
}

// Consumer: listen for new orders
func ConsumeOrders(handler func(OrderMessage)) {
	InitBrokers()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "order-service",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	go func() {
		defer r.Close()
		log.Println("[Kafka] Consumer started...")
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[Kafka] Consumer error: %v", err)
				continue
			}

			var msg OrderMessage
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("[Kafka] JSON decode error: %v", err)
				continue
			}

			log.Printf("[Kafka] Consumed order %s", msg.OrderID)
			handler(msg)
		}
	}()
}
