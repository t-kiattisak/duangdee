package main

import (
	"context"
	"log"
	"os"

	"duangdee/pkg/logger"
	"duangdee/pkg/middleware"
	handler "duangdee/tarot/internal/delivery/http"
	"duangdee/tarot/internal/repository"
	"duangdee/tarot/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	sysLogger := logger.New("tarot-service")

	// 1. Database Connection Pool Setup (PostgreSQL tarot_db)
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgres://postgres:secretpassword@postgres:5432/tarot_db?sslmode=disable"
	}

	dbPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("Connected to PostgreSQL (tarot_db) successfully")

	// 2. Resolve Local Assets Directory Path
	assetsDir := os.Getenv("ASSETS_DIR")
	if assetsDir == "" {
		assetsDir = "./assets/cards"
	}

	// 3. Clean Architecture Wiring
	cardRepo := repository.NewCardRepository(dbPool)
	cardUsecase := usecase.NewCardUsecase(cardRepo)
	cardHandler := handler.NewCardHandler(cardUsecase, assetsDir)

	// 4. Go Fiber App Initialization
	app := fiber.New(fiber.Config{
		AppName: "Duangdee Tarot Service",
	})

	app.Use(cors.New())
	app.Use(middleware.Logger(sysLogger))

	// Health Check Route
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "tarot-service",
		})
	})

	// Dedicated Secure Static Image Streaming API
	app.Get("/assets/cards/:filename", cardHandler.ServeImage)

	// API V1 Routes
	api := app.Group("/api/v1/tarot")
	api.Get("/cards", cardHandler.ListCards)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Tarot Core Service running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
