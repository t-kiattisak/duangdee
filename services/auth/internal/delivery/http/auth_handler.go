package http

import (
	"errors"
	"log"

	"duangdee/auth/internal/domain"
	"duangdee/auth/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Register handles HTTP POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body payload",
		})
	}

	// Payload validation
	if req.Username == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username, password, first_name, and last_name are required fields",
		})
	}

	if len(req.Username) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "username must be at least 3 characters long",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password must be at least 6 characters long",
		})
	}

	// Execute Business Usecase
	res, err := h.authUsecase.Register(c.Context(), &req)
	if err != nil {
		log.Printf("Registration Error: %v", err)
		if errors.Is(err, usecase.ErrUsernameAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "username is already taken",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to complete registration",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "registration successful",
		"data":    res,
	})
}
