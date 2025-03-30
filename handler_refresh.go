package main

import (
	"net/http"
	"strings"
	"time"

	"chirpy/internal/auth"
)

// Add this type for the response
type RefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		RespondWithError(w, http.StatusUnauthorized, "Missing authorization header", nil)
		return
	}

	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")
	if refreshToken == authHeader {
		RespondWithError(w, http.StatusUnauthorized, "Invalid authorization header format", nil)
		return
	}

	// Get refresh token from database
	token, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	// Check if the token is expired
	if time.Now().After(token.ExpiresAt) {
		RespondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	// Check if the token is revoked
	if token.RevokedAt.Valid {
		RespondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}

	// Create new access token
	accessToken, err := auth.MakeJWT(
		token.UserID,
		cfg.tokenSecret,
		time.Hour,
	)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to create access token", err)
		return
	}

	response := RefreshResponse{
		Token: accessToken,
	}

	RespondWithJSON(w, http.StatusOK, response)
}
