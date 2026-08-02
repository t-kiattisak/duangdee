package main

import (
	"os"

	"duangdee/pkg/logger"
	"duangdee/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Duangdee Payment Service",
	})

	sysLogger := logger.New("payment-service")

	app.Use(cors.New())
	app.Use(middleware.Logger(sysLogger))

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "payment-service",
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
