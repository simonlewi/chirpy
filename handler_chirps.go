package main

import "net/http"

func (cfg *apiConfig) ChirpsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cfg.GetChirps(w, r)
	case http.MethodPost:
		cfg.CreateChirp(w, r)
	case http.MethodDelete:
		cfg.DeleteChirpHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
