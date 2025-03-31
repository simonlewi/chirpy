package main

import (
	"chirpy/internal/auth"
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
	tokenSecret    string
	tokenService   *auth.TokenService
	apiKey         string
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
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
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
		dbQueries:    database.New(dbConn),
		platform:     platform,
		tokenSecret:  jwtSecret,
		tokenService: auth.NewTokenService(database.New(dbConn), jwtSecret),
		apiKey:       polkaKey,
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
	mux.HandleFunc("/api/refresh", cfg.RefreshHandler)
	mux.HandleFunc("/api/revoke", cfg.RevokeHandler)

	mux.HandleFunc("/api/polka/webhooks", cfg.HandlerWebhooks)

	mux.HandleFunc("/admin/metrics", cfg.MetricsHandler)
	mux.HandleFunc("/admin/reset", cfg.ResetHandler)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())

}
