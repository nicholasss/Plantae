package main

import (
	"database/sql"
	"encoding/json"
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

// response types

type createUserRequest struct {
	CreatedBy   string `json:"createdBy"`
	UpdatedBy   string `json:"updatedBy"`
	IsAdmin     bool   `json:"isAdmin"`
	Email       string `json:"email"`
	RawPassword string `json:"rawPassword"`
}

// === Middleware Functions ===

func (cfg *apiConfig) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// === Utility Response Handlers ===

func respondWithError(err error, w http.ResponseWriter) {
}

// === Handler Functions ===

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest createUserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createUserRequest)
	if err != nil {
		// respond with error
	}

	// check request params
	if createUserRequest.Email == "" {
		// respond with error
	}
	if createUserRequest.RawPassword == "" {
		// respond with error
	}
	if createUserRequest.CreatedBy == "" {
		// respond with error
	}
	if createUserRequest.UpdatedBy == "" {
		// respond with error
	}
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

	// health endpoint
	mux.Handle("GET /health", cfg.logMW(http.HandlerFunc(healthHandler)))

	// user endpoints
	mux.Handle("POST /api/v1/createuser", cfg.logMW(http.HandlerFunc(cfg.createUserHandler)))

	log.Printf("Server is now online.\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
