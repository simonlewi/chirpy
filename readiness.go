package main

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// If it's not, respond with a 405 (Method Not Allowed)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
