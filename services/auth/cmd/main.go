package main

import (
	"context"
	"os"

	"duangdee/pkg/kafka"
	"duangdee/pkg/logger"
	"duangdee/pkg/middleware"

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

	// Attach HTTP Request/Response Metadata Logger Middleware
	app.Use(middleware.Logger(sysLogger))

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

		return c.JSON(fiber.Map{
			"status":  "published_to_kafka",
			"message": "Test event published successfully",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sysLogger.Info("Service started", map[string]interface{}{"port": port})
	if err := app.Listen(":" + port); err != nil {
		sysLogger.Error(err, "Failed to start server", nil)
		os.Exit(1)
	}
}
