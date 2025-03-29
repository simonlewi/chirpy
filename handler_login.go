package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
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

	// Determine token expiration time
	const maxExpirationTime = 60 * 60 // 1 hour in seconds
	expirationTime := maxExpirationTime
	if params.ExpiresInSeconds > 0 {
		if params.ExpiresInSeconds > maxExpirationTime {
			expirationTime = maxExpirationTime
		} else {
			expirationTime = params.ExpiresInSeconds
		}
	}

	// Create JWT token
	token, err := auth.MakeJWT(
		dbUser.ID,
		cfg.secret,
		time.Duration(expirationTime)*time.Second,
	)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Create response
	response := LoginResponse{
		User: User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		},
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Couldn't encode response", http.StatusInternalServerError)
		return
	}
}
