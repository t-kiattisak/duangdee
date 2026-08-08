package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"duangdee/auth/internal/domain"
	"duangdee/auth/pkg/jwt"
	"duangdee/pkg/kafka"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameAlreadyExists = errors.New("username is already taken")
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrAccountInactive       = errors.New("user account is inactive")
)

type authUsecase struct {
	userRepo   domain.UserRepository
	tokenMaker *jwt.TokenMaker
	producer   *kafka.Producer
}

// NewAuthUsecase returns a new instance of domain.AuthUsecase
func NewAuthUsecase(repo domain.UserRepository, tokenMaker *jwt.TokenMaker, producer *kafka.Producer) domain.AuthUsecase {
	return &authUsecase{
		userRepo:   repo,
		tokenMaker: tokenMaker,
		producer:   producer,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 1. Check if username already exists
	exists, err := u.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing username: %w", err)
	}
	if exists {
		return nil, ErrUsernameAlreadyExists
	}

	// 2. Hash Password using Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Optional: Parse Birth Date & Calculate Zodiac Sign
	var birthDate *time.Time
	var zodiacSign *string

	if req.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err == nil {
			birthDate = &parsedDate
			sign := computeZodiacSign(parsedDate.Month(), parsedDate.Day())
			zodiacSign = &sign
		}
	}

	// 4. Create User domain entity
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BirthDate:    birthDate,
		ZodiacSign:   zodiacSign,
		IsActive:     true,
	}

	// 5. Persist to PostgreSQL database
	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// 6. Issue JWT Tokens (Access & Refresh)
	accessToken, refreshToken, err := u.tokenMaker.GenerateTokens(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// 7. Publish `user.registered` event to Kafka (Async/Non-blocking error)
	if u.producer != nil {
		_ = u.producer.Publish(ctx, "user.registered", user.ID, kafka.EventMessage{
			EventID:   uuid.New().String(),
			EventType: "user.registered",
			Data: map[string]interface{}{
				"user_id":    user.ID,
				"username":   user.Username,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
			},
		})
	}

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// 1. Fetch user by username
	user, err := u.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verify password with Bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// 3. Verify user account status
	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	// 4. Generate new JWT Access Token & Refresh Token
	accessToken, refreshToken, err := u.tokenMaker.GenerateTokens(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &domain.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Helper: Calculate Western Zodiac Sign from Month and Day
func computeZodiacSign(month time.Month, day int) string {
	switch month {
	case time.January:
		if day <= 19 {
			return "Capricorn"
		}
		return "Aquarius"
	case time.February:
		if day <= 18 {
			return "Aquarius"
		}
		return "Pisces"
	case time.March:
		if day <= 20 {
			return "Pisces"
		}
		return "Aries"
	case time.April:
		if day <= 19 {
			return "Aries"
		}
		return "Taurus"
	case time.May:
		if day <= 20 {
			return "Taurus"
		}
		return "Gemini"
	case time.June:
		if day <= 20 {
			return "Gemini"
		}
		return "Cancer"
	case time.July:
		if day <= 22 {
			return "Cancer"
		}
		return "Leo"
	case time.August:
		if day <= 22 {
			return "Leo"
		}
		return "Virgo"
	case time.September:
		if day <= 22 {
			return "Virgo"
		}
		return "Libra"
	case time.October:
		if day <= 22 {
			return "Libra"
		}
		return "Scorpio"
	case time.November:
		if day <= 21 {
			return "Scorpio"
		}
		return "Sagittarius"
	case time.December:
		if day <= 21 {
			return "Sagittarius"
		}
		return "Capricorn"
	}
	return "Unknown"
}
