package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ProfaneFlag is a function that replaces profane words with asterisks
func ProfaneFlag(inputText string) string {

	profaneWordMap := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	words := strings.Split(inputText, " ")
	for i, word := range words {
		if profaneWordMap[strings.ToLower(word)] {
			words[i] = "****"
		}
	}
	cleanedText := strings.Join(words, " ")
	return cleanedText
}

// HandlerProfane is an HTTP handler that cleans profanity from a chirp
func HandlerProfane(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		return
	}

	var req struct {
		Body string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if len(req.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Clean the text using profanity filter
	cleanedBody := ProfaneFlag(req.Body)

	RespondWithJSON(w, http.StatusOK, map[string]string{
		"cleaned_body": cleanedBody,
	})
}
