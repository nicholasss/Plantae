package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicholasss/plantae/internal/database"
)

// === Global Types ===

type apiConfig struct {
	db        *database.Queries
	localAddr string
	port      string
}

// === Utilities Response Types ===

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

// === Utility Functions ===

func loadApiConfig() (*apiConfig, error) {
	// setting up connection to the database
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	dbURL := os.Getenv("GOOSE_DBSTRING")
	if dbURL == "" {
		return nil, err
	}
	log.Print("Connected to database succesfully.")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	dbQueries := database.New(db)

	cfg := &apiConfig{
		db:        dbQueries,
		localAddr: os.Getenv("LOCAL_ADDRESS"),
		port:      ":" + os.Getenv("PORT"),
	}

	return cfg, nil
}

// === Utility Response Handlers ===

func respondWithError(error error, code int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if error != nil {
		errorString := error.Error()
		errorResponse := errorResponse{ErrorMessage: errorString}
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
