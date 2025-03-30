package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) CreateChirp(w http.ResponseWriter, r *http.Request) {
	// Get the bearer token from the Authorization header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Validate chirp length
	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Clean the text using profanity filter
	cleanedBody := ProfaneFlag(params.Body)

	// Create a new chirp using userID from JWT
	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// Convert database.Chirp to your API Chirp
	responseChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	RespondWithJSON(w, http.StatusCreated, responseChirp)
}

func (cfg *apiConfig) GetChirps(w http.ResponseWriter, r *http.Request) {
	// Call the SQLC-generated funtion to get ALL users
	dbUsers, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %v", err)
		http.Error(w, "Couldn't get chirps", http.StatusInternalServerError)
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbUsers {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(chirps)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Couldn't encode response", http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) GetChirpID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			RespondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get chirp", err)
		return
	}

	responseChirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	RespondWithJSON(w, http.StatusOK, responseChirp)
}

func (cfg *apiConfig) ChirpsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cfg.GetChirps(w, r)
		return
	} else if r.Method == http.MethodPost {
		cfg.CreateChirp(w, r)
		return
	} else {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}
}
