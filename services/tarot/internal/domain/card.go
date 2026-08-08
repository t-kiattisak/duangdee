package domain

import (
	"context"
)

// TarotCard represents a tarot card catalog entity
type TarotCard struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	ArcanaType string  `json:"arcana_type"`
	Suit       *string `json:"suit,omitempty"`
	Number     int     `json:"number"`
	Element    *string `json:"element,omitempty"`
	ImageURL   string  `json:"image_url"`
}

// CardRepository defines the persistence contract for Tarot Cards
type CardRepository interface {
	GetAllCards(ctx context.Context, arcanaType string) ([]*TarotCard, error)
	GetCardByID(ctx context.Context, id int) (*TarotCard, error)
}

// CardUsecase defines the business logic contract for Tarot Cards
type CardUsecase interface {
	ListCards(ctx context.Context, arcanaType string) ([]*TarotCard, error)
	GetCardDetails(ctx context.Context, id int) (*TarotCard, error)
}
