package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

// returns true if the platform is production
func platformProduction(cfg *apiConfig) (bool, error) {
	if cfg.platform == "" {
		return false, errors.New("platform variable is missing")
	}

	if cfg.platform == "production" {
		return true, nil
	}

	return false, nil
}

// returns true if the platform is not production
func platformNotProduction(cfg *apiConfig) (bool, error) {
	if cfg.platform == "" {
		return false, errors.New("platform variable is missing")
	}

	if cfg.platform == "production" {
		return false, nil
	}

	return true, nil
}

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

	// checking the config
	if cfg.localAddr == "" {
		log.Print("WARNING: 'LOCAL_ADDRESS' is empty, please check .env")
	}
	if cfg.platform == "" {
		log.Print("WARNING: 'PLATFORM' is empty, please check .env")
	}
	if cfg.port == "" {
		log.Print("WARNING: 'PORT' is empty, please check .env")
	}
	if cfg.JWTSecret == "" {
		log.Print("WARNING: 'JWT_SECRET' is empty, please check .env")
	}
	if cfg.superAdminToken == "" {
		log.Print("WARNING: 'SUPER_ADMIN_TOKEN' is empty, please check .env")
	}

	log.Printf("Platform loaded as %q.", cfg.platform)

	return cfg, nil
}

// === Utility Response Handlers ===

// TODO: function to respond due to a wrong platform for action
// some kind of enum for action? reset, promote, etc.

// TODO: send out a more generic error to client
func respondWithError(err error, code int, w http.ResponseWriter) {
	log.Printf("ERROR: Sending error to client: %q", err)

	// response side
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err != nil {
		errorString := err.Error()
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
