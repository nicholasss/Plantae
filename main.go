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
	db        *database.Queries
	localAddr string
	port      string
}

// === Handler Functions ===

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
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
	if dbURL == "" {
		log.Fatalf("Unable to find database string with: %q", "GOOSE_DBSTRING")
	}
	log.Print("Connected to database succesfully.")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Unable to open connection to database: %s", err)
	}

	dbQueries := database.New(db)

	cfg := &apiConfig{
		db:        dbQueries,
		localAddr: os.Getenv("LOCAL_ADDRESS"),
		port:      ":" + os.Getenv("PORT"),
	}

	mux := http.NewServeMux()

	// health endpoint
	mux.Handle("GET /health", cfg.logMW(http.HandlerFunc(healthHandler)))

	// user endpoints
	mux.Handle("POST /api/v1/createuser", cfg.logMW(http.HandlerFunc(cfg.createUserHandler)))

	log.Printf("Server is now online at http://%s%s.\n", cfg.localAddr, cfg.port)
	log.Fatal(http.ListenAndServe(cfg.port, mux))
}
