package main

import (
	"chirpy/internal/database"
	"database/sql"
	"fmt"
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
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	cfg := &apiConfig{
		dbQueries: database.New(db),
	}

	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", cfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("/api/healthz", HandlerReadiness)
	mux.HandleFunc("/api/validate_chirp", HandlerProfane)
	mux.HandleFunc("/admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("/admin/reset", cfg.ResetHandler)

	httpServer := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// If it's not GET, respond with a 405 (Method Not Allowed)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	html, err := os.ReadFile("admin.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, string(html), cfg.fileserverHits.Load())
}
