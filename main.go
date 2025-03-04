package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicholasss/plantdata/internal/database"

	_ "github.com/lib/pq"
)

// === Global Types ===

type apiConfig struct {
	db *database.Queries
}

// === Middleware Functions ===

func (cfg *apiConfig) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// === Handler Functions ===

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// === Main Function ===

func main() {
	log.Printf("Staring server\n")

	// setting up connection to the database
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Unable to load '.env'.\n")
	}

	dbURL := os.Getenv("GOOSE_DBSTRING")
	log.Printf("Database URL: %s\n", dbURL)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to open connection to database: %s", err)
	}

	dbQueries := database.New(db)

	cfg := &apiConfig{
		db: dbQueries,
	}

	mux := http.NewServeMux()

	mux.Handle("GET /health", cfg.logMW(http.HandlerFunc(healthHandler)))

	log.Printf("Server is now online.\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
