package usecase

import (
	"context"

	"duangdee/tarot/internal/domain"
)

type cardUsecase struct {
	cardRepo domain.CardRepository
}

// NewCardUsecase returns a new instance of domain.CardUsecase
func NewCardUsecase(cardRepo domain.CardRepository) domain.CardUsecase {
	return &cardUsecase{cardRepo: cardRepo}
}

func (u *cardUsecase) ListCards(ctx context.Context, arcanaType string) ([]*domain.TarotCard, error) {
	return u.cardRepo.GetAllCards(ctx, arcanaType)
}

func (u *cardUsecase) GetCardDetails(ctx context.Context, id int) (*domain.TarotCard, error) {
	return u.cardRepo.GetCardByID(ctx, id)
}
