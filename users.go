package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params createUserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Attempting to create user with email: %s", params.Email)

	// Has the password before storing
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	createUserParams := database.CreateUserParams{
		Email:   params.Email,
		Column2: hashedPassword, // Store the hashed password
	}

	// Call the SQLC-generated function to create a user
	dbUser, err := cfg.dbQueries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, fmt.Sprintf("Couldn't create user: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert database.User to your API User (without password)
	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(responseUser)
	if err != nil {
		http.Error(w, "Couldn't encode user", http.StatusInternalServerError)
		return
	}

}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {
	// Call the SQLC-generated funtion to get ALL users
	dbUsers, err := cfg.dbQueries.GetUsers(r.Context())
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Couldn't get users", http.StatusInternalServerError)
		return
	}

	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, User{
			ID:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Couldn't encode response", http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cfg.createUser(w, r)
		return
	}
	if r.Method == http.MethodGet {
		cfg.getUsers(w, r)
		return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}
