package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

// Chirp represents a message in the system
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
	// Get the author_id from query parameters
	authorIDStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	var dbChirps []database.Chirp
	var err error

	if authorIDStr != "" {
		// Parse and validate the UUID
		var authorID uuid.UUID
		authorID, err = uuid.Parse(authorIDStr)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}

		// Get chirps filtered by author
		dbChirps, err = cfg.dbQueries.GetChirpsByAuthor(r.Context(), authorID)
	} else {
		// Get all chirps if no author filter
		dbChirps, err = cfg.dbQueries.GetChirps(r.Context())
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	// Convert database chirps to response format
	chirps := make([]Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		if dbChirp.ID == uuid.Nil || dbChirp.Body == "" {
			continue
		}

		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	RespondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) GetChirpID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

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

	case http.MethodDelete:
		cfg.DeleteChirpHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
