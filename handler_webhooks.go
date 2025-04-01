package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

// WebhookHandler Documentation
//
// Endpoint: POST /api/polka/webhooks
//
// Description:
// Handles webhook events for user upgrades to Chirpy Red status. This endpoint
// processes notifications from the payment platform to upgrade user memberships.
//
// Request Body:
// {
//   "event": "user.upgraded",
//   "data": {
//     "user_id": "uuid-string"
//   }
// }
//
// Response Status Codes:
// - 204: Success (No Content)
//       Returned when:
//       * Event processed successfully
//       * Event type is not "user.upgraded" (silent ignore)
// - 400: Bad Request
//       Returned when:
//       * Invalid JSON payload
//       * Invalid UUID format
// - 404: Not Found
//       Returned when:
//       * User ID cannot be found in database
// - 405: Method Not Allowed
//       Returned when:
//       * HTTP method is not POST
//
// Example Usage:
// curl -X POST http://localhost:8080/api/polka/webhooks \
//   -H "Content-Type: application/json" \
//   -d '{"event":"user.upgraded","data":{"user_id":"3311741c-680c-4546-99f3-fc9efac2036c"

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
