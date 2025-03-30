package main

import "net/http"

func (cfg *apiConfig) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cfg.createUser(w, r)
	case http.MethodGet:
		cfg.getUsers(w, r)
	case http.MethodPut:
		cfg.UpdateUserHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
