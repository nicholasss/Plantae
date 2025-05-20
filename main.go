package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/nicholasss/plantdata/internal/auth"
	"github.com/nicholasss/plantdata/internal/database"

	_ "github.com/lib/pq"
)

// === Global Types ===

type apiConfig struct {
	db *database.Queries
}

// request types

type createUserRequest struct {
	CreatedBy   string `json:"createdBy"`
	UpdatedBy   string `json:"updatedBy"`
	IsAdmin     bool   `json:"isAdmin"`
	Email       string `json:"email"`
	RawPassword string `json:"rawPassword"`
}

// response types

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

// === Middleware Functions ===

func (cfg *apiConfig) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
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
		respondWithError(nil, http.StatusBadRequest, w)
		return
	}
	if createUserRequest.RawPassword == "" { // may not need to check due to sql.NullString type
		respondWithError(nil, http.StatusBadRequest, w)
		return
	}
	if createUserRequest.CreatedBy == "" {
		respondWithError(nil, http.StatusBadRequest, w)
		return
	}
	if createUserRequest.UpdatedBy == "" {
		respondWithError(nil, http.StatusBadRequest, w)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(createUserRequest.RawPassword)
	createUserRequest.RawPassword = "" // GC collection
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}
	validHashedPassword := sql.NullString{
		String: hashedPassword,
		Valid:  true,
	}

	// CreateUserParams struct
	createUserParams := database.CreateUserParams{
		CreatedBy:      createUserRequest.CreatedBy,
		UpdatedBy:      createUserRequest.UpdatedBy,
		IsAdmin:        createUserRequest.IsAdmin,
		Email:          createUserRequest.Email,
		HashedPassword: validHashedPassword,
	}

	// add user to database
	userRecord, err := cfg.db.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// return the userRecord without password
	userData, err := json.Marshal(userRecord)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userData)
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
