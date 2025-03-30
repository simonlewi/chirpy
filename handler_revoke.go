package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
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

	// Add the token to your revoked tokens list or database
	err := cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
