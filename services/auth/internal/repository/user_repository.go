package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"duangdee/auth/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository returns a new instance of domain.UserRepository
func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (id, username, password_hash, first_name, last_name, avatar_url, birth_date, zodiac_sign, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		u.ID, u.Username, u.PasswordHash, u.FirstName, u.LastName,
		u.AvatarURL, u.BirthDate, u.ZodiacSign, u.IsActive, u.CreatedAt, u.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert user into db: %w", err)
	}
	return nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, password_hash, first_name, last_name, avatar_url, birth_date, zodiac_sign, is_active, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.FirstName, &u.LastName,
		&u.AvatarURL, &u.BirthDate, &u.ZodiacSign, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Return nil, nil when user not found
		}
		return nil, fmt.Errorf("error querying user by username: %w", err)
	}
	return u, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
	var exists bool
	err := r.db.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking username existence: %w", err)
	}
	return exists, nil
}
