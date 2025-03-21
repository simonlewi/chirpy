package main

import (
	"encoding/json"
	"log"
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

// respondWithError is a helper function to respond with an error message
func respondWithError(w http.ResponseWriter, code int, msg string) {

	respondWithJSON(w, code, map[string]string{"error": msg})
}

// respondWithJSON is a helper function to respond with a JSON payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// If encoding fails, log the error and send a plain text response
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

// HandlerProfane is an HTTP handler that cleans profanity from a chirp
func HandlerProfane(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var req struct {
		Body string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Clean the text using our profanity filter
	cleanedBody := ProfaneFlag(req.Body)

	respondWithJSON(w, http.StatusOK, map[string]string{
		"cleaned_body": cleanedBody,
	})
}
