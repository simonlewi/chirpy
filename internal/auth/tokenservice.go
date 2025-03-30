package auth

import (
	"chirpy/internal/database"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TokenService struct {
	db          *database.Queries
	tokenSecret string
}

func NewTokenService(db *database.Queries, secret string) *TokenService {
	return &TokenService{
		db:          db,
		tokenSecret: secret,
	}
}

func (s *TokenService) CreateTokenPair(userID uuid.UUID, expiresInSeconds int) (string, string, error) {
	// Create access token
	accessToken, err := MakeJWT(
		userID,
		s.tokenSecret,
		time.Hour,
	)
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	refreshToken, err := MakeRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Store refresh token in database
	refreshTokenExpiry := time.Now().Add(60 * 24 * time.Hour) // 60 days storage
	err = s.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		ExpiresAt: refreshTokenExpiry,
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Parse and validate the JWT token
func (s *TokenService) ValidateToken(token string) (uuid.UUID, error) {
	// Parse and validate the JWT token
	userID, err := ValidateJWT(token, s.tokenSecret)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

// Check if token exists and is revoked in the database
func (s *TokenService) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	refreshToken, err := s.db.GetRefreshToken(ctx, token)
	if err != nil {
		return true, err
	}

	// Check if token is revoked (RevokedAt has a non-zero value)
	return refreshToken.RevokedAt.Valid, nil
}

func (s *TokenService) RevokeToken(ctx context.Context, token string) error {
	return s.db.RevokeRefreshToken(ctx, token)
}
