package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}
type Order struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	UserID     uint   `json:"user_id"`
	Status     string `json:"status"`
	TotalCents int64  `json:"total_cents"`
}

type Job struct{ OrderID uint }

var jobs chan Job

func worker(db *gorm.DB, id int) {
	for j := range jobs {
		time.Sleep(time.Second)
		status := "failed"
		if rand.Intn(100) < 70 {
			status = "confirmed"
		}
		db.Model(&Order{}).Where("id=? AND status='pending'", j.OrderID).
			Update("status", status)
		fmt.Printf("Worker %d -> order %d: %s\n", id, j.OrderID, status)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	workers := 3
	if w, _ := strconv.Atoi(os.Getenv("WORKERS")); w > 0 {
		workers = w
	}
	jobs = make(chan Job, 16)

	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&User{}, &Order{})

	for i := 1; i <= workers; i++ {
		go worker(db, i)
	}

	app := fiber.New()

	app.Post("/users", func(c *fiber.Ctx) error {
		var u User
		if c.BodyParser(&u) != nil || db.Create(&u).Error != nil {
			return fiber.ErrBadRequest
		}
		return c.JSON(u)
	})
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		var u User
		if db.First(&u, c.Params("id")).Error != nil {
			return fiber.ErrNotFound
		}
		return c.JSON(u)
	})
	app.Post("/orders", func(c *fiber.Ctx) error {
		var o Order
		if c.BodyParser(&o) != nil {
			return fiber.ErrBadRequest
		}
		o.Status = "pending"
		db.Create(&o)
		return c.JSON(o)
	})
	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		var o Order
		if db.First(&o, c.Params("id")).Error != nil {
			return fiber.ErrNotFound
		}
		return c.JSON(o)
	})
	app.Post("/orders/:id/confirm", func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		var o Order
		if db.First(&o, id).Error != nil {
			return fiber.ErrNotFound
		}
		select {
		case jobs <- Job{OrderID: o.ID}:
			return c.JSON(fiber.Map{"message": "queued"})
		default:
			return fiber.ErrServiceUnavailable
		}
	})

	go app.Listen(":8080")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	fmt.Println("bye!")
	close(jobs)
}
