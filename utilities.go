package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nicholasss/plantae/internal/database"
)

// === Global Types ===

type apiConfig struct {
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	db                   *database.Queries
	localAddr            string
	platform             string
	port                 string
	JWTSecret            string
	superAdminToken      string
}

// === Utilities Response Types ===

type errorResponse struct {
	Error string `json:"error"`
}

// === Utility Functions ===

func loadApiConfig() (*apiConfig, error) {
	// loading vars from .env
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	// connecting to database
	dbURL := os.Getenv("GOOSE_DBSTRING")
	if dbURL == "" {
		return nil, err
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	dbQueries := database.New(db)
	log.Print("Connected to database succesfully.")

	// additional vars, configuration, and return

	cfg := &apiConfig{
		accessTokenDuration:  time.Hour * 2,
		refreshTokenDuration: time.Hour * 24 * 30,
		db:                   dbQueries,
		localAddr:            os.Getenv("LOCAL_ADDRESS"),
		platform:             os.Getenv("PLATFORM"),
		port:                 ":" + os.Getenv("PORT"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		superAdminToken:      os.Getenv("SUPER_ADMIN_TOKEN"),
	}

	return cfg, nil
}

// === Utility Response Handlers ===

func respondWithError(error error, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if error != nil {
		errorString := error.Error()
		errorResponse := errorResponse{Error: errorString}
		errorData, err := json.Marshal(errorResponse)
		if err != nil {
			log.Printf("Error occured marshaling error response: %q", err)
			return
		}

		w.Write(errorData)
		return
	}

	w.Write([]byte(`{"error":"internal error"}`))
}
