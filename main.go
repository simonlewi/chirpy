package main

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	secret         string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("PLATFORM must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}
	defer dbConn.Close()

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	cfg := &apiConfig{
		dbQueries: database.New(dbConn),
		platform:  platform,
		secret:    jwtSecret,
	}

	mux := http.NewServeMux()

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(fileHandler))

	mux.HandleFunc("/api/healthz", HandlerReadiness)
	mux.HandleFunc("/api/validate_chirp", HandlerProfane)
	mux.HandleFunc("/api/chirps", cfg.ChirpsHandler)
	mux.HandleFunc("/api/chirps/{chirpID}", cfg.GetChirpID)

	mux.HandleFunc("/api/users", cfg.UsersHandler)
	mux.HandleFunc("/api/login", cfg.LoginHandler)

	mux.HandleFunc("/admin/metrics", cfg.MetricsHandler)
	mux.HandleFunc("/admin/reset", cfg.ResetHandler)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())

}
