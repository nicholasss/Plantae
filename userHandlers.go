package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// request types

type createUserRequest struct {
	CreatedBy   string `json:"createdBy"`
	UpdatedBy   string `json:"updatedBy"`
	Email       string `json:"email"`
	RawPassword string `json:"rawPassword"`
}

// === User Handler Functions ===

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
