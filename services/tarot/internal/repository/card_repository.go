package repository

import (
	"context"
	"errors"
	"fmt"

	"duangdee/tarot/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type cardRepository struct {
	db *pgxpool.Pool
}

// NewCardRepository returns a new instance of domain.CardRepository
func NewCardRepository(db *pgxpool.Pool) domain.CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) GetAllCards(ctx context.Context, arcanaType string) ([]*domain.TarotCard, error) {
	query := `
		SELECT id, name, arcana_type, suit, number, element, image_url
		FROM tarot_cards
	`
	var args []interface{}
	if arcanaType != "" {
		query += " WHERE arcana_type = $1"
		args = append(args, arcanaType)
	}
	query += " ORDER BY id ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cards: %w", err)
	}
	defer rows.Close()

	var cards []*domain.TarotCard
	for rows.Next() {
		card := &domain.TarotCard{}
		err := rows.Scan(&card.ID, &card.Name, &card.ArcanaType, &card.Suit, &card.Number, &card.Element, &card.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("error scanning card row: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (r *cardRepository) GetCardByID(ctx context.Context, id int) (*domain.TarotCard, error) {
	query := `
		SELECT id, name, arcana_type, suit, number, element, image_url
		FROM tarot_cards
		WHERE id = $1
	`
	card := &domain.TarotCard{}
	err := r.db.QueryRow(ctx, query, id).Scan(&card.ID, &card.Name, &card.ArcanaType, &card.Suit, &card.Number, &card.Element, &card.ImageURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying card by id: %w", err)
	}

	return card, nil
}
