package main

import (
	"context"
	"log"
	"os"

	"duangdee/pkg/kafka"
	"duangdee/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Duangdee Auth Service",
	})

	// Initialize Shared Zerolog JSON Logger
	sysLogger := logger.New("auth-service")

	app.Use(cors.New())

	// Middleware to log HTTP Requests in Structured JSON for Kibana
	app.Use(func(c *fiber.Ctx) error {
		sysLogger.Info("Incoming HTTP Request", map[string]interface{}{
			"method": c.Method(),
			"path":   c.Path(),
			"ip":     c.IP(),
		})
		return c.Next()
	})

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	// Test Producer for Kafka
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}
	producer := kafka.NewProducer([]string{kafkaBrokers})
	defer producer.Close()

	// Test Route: Emit event to Kafka
	app.Post("/api/v1/auth/test-event", func(c *fiber.Ctx) error {
		err := producer.Publish(context.Background(), "user.registered", "test_user_1", kafka.EventMessage{
			EventID:   "evt_test_123",
			EventType: "user.registered",
			Data: map[string]interface{}{
				"user_id": "usr_test_123",
				"email":   "test@duangdee.com",
			},
		})
		if err != nil {
			sysLogger.Error(err, "Failed to publish test event to Kafka", nil)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		sysLogger.Info("Published user.registered event to Kafka successfully", map[string]interface{}{
			"topic": "user.registered",
		})

		return c.JSON(fiber.Map{"status": "published_to_kafka"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Auth Service running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
