package main

import (
	"context"
	"log"
	"os"

	"duangdee/auth/internal/delivery/http"
	"duangdee/auth/internal/repository"
	"duangdee/auth/internal/usecase"
	"duangdee/auth/pkg/jwt"
	"duangdee/pkg/kafka"
	"duangdee/pkg/logger"
	"duangdee/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	sysLogger := logger.New("auth-service")

	// 1. Database Connection Pool Setup (PostgreSQL auth_db)
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgres://postgres:secretpassword@postgres:5432/auth_db?sslmode=disable"
	}

	dbPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("Connected to PostgreSQL (auth_db) successfully")

	// 2. Kafka Producer Setup
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "kafka:9092"
	}
	producer := kafka.NewProducer([]string{kafkaBrokers})
	defer producer.Close()

	// 3. JWT Token Maker Setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretduangdeeauthkey123"
	}
	tokenMaker := jwt.NewTokenMaker(jwtSecret)

	// 4. Clean Architecture Dependency Wiring
	userRepo := repository.NewUserRepository(dbPool)
	authUsecase := usecase.NewAuthUsecase(userRepo, tokenMaker, producer)
	authHandler := http.NewAuthHandler(authUsecase)

	// 5. Go Fiber Web Framework Initialization
	app := fiber.New(fiber.Config{
		AppName: "Duangdee Auth Service",
	})

	app.Use(cors.New())
	app.Use(middleware.Logger(sysLogger))

	// Health Check Route
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "auth-service",
		})
	})

	// API V1 Routes
	api := app.Group("/api/v1/auth")
	api.Post("/register", authHandler.Register)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Auth Service running on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
