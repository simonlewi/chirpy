package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// Endpoint: POST /api/polka/webhooks

type WebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) HandlerWebhooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	expectedAPIKey := os.Getenv("POLKA_KEY")
	if apiKey != expectedAPIKey {
		RespondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	webhookReq := WebhookRequest{}
	err = decoder.Decode(&webhookReq)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if webhookReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(webhookReq.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	_, err = cfg.dbQueries.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "User not found", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
