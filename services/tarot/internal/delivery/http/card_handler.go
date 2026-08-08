package http

import (
	"log"
	"os"
	"path/filepath"

	"duangdee/tarot/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type CardHandler struct {
	cardUsecase domain.CardUsecase
	assetsDir   string
}

func NewCardHandler(cardUsecase domain.CardUsecase, assetsDir string) *CardHandler {
	return &CardHandler{
		cardUsecase: cardUsecase,
		assetsDir:   assetsDir,
	}
}

// ListCards handles HTTP GET /api/v1/tarot/cards
func (h *CardHandler) ListCards(c *fiber.Ctx) error {
	arcanaType := c.Query("arcana")
	cards, err := h.cardUsecase.ListCards(c.Context(), arcanaType)
	if err != nil {
		log.Printf("Error fetching cards: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to retrieve cards catalog",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"total": len(cards),
			"cards": cards,
		},
	})
}

// ServeImage handles secure dedicated image stream GET /assets/cards/:filename
func (h *CardHandler) ServeImage(c *fiber.Ctx) error {
	rawFilename := c.Params("filename")

	// 1. Strict Sanitization: Use filepath.Base to neutralize Path Traversal Attacks (e.g. "../../etc/passwd")
	filename := filepath.Base(rawFilename)

	// 2. Strict Extension Whitelist Enforcement
	ext := filepath.Ext(filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid or forbidden file extension",
		})
	}

	// 3. Resolve absolute file path safely
	targetPath := filepath.Join(h.assetsDir, filename)

	// Check file existence
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "card image not found",
		})
	}

	// 4. Set High Performance HTTP Caching Header (1 Day = 86400s)
	c.Set("Cache-Control", "public, max-age=86400")

	// 5. Safely stream local file to response
	return c.SendFile(targetPath)
}
