package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var params LoginRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding login request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get user from database by email
	dbUser, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		// If the user is not found, we can return a generic error message
		// to avoid leaking information about whether the email exists
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := auth.CheckPasswordHash(dbUser.HashedPassword, params.Password); err != nil {
		http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
		return
	}

	// Create tokens using TokenService
	accessToken, refreshToken, err := cfg.tokenService.CreateTokenPair(dbUser.ID, params.ExpiresInSeconds)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error creating tokens", err)
		return
	}

	// Create response
	response := LoginResponse{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		IsChirpyRed:  dbUser.IsChirpyRed.Bool,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	RespondWithJSON(w, http.StatusOK, response)

}
