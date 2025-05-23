package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// request types

type AdminStatusRequest struct {
	ID uuid.UUID `json:"id"`
}

type CreateUserRequest struct {
	CreatedBy   string `json:"createdBy"`
	UpdatedBy   string `json:"updatedBy"`
	Email       string `json:"email"`
	RawPassword string `json:"rawPassword"`
}

// login endpoint
type UserLoginRequest struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}
type UserLoginResponse struct {
	ID           uuid.UUID `json:"id"`
	IsAdmin      bool      `json:"is_admin"`
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

// === User Handler Functions ===

func (cfg *apiConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	// ensure development platform
	if cfg.platform == "production" || cfg.platform == "" {
		log.Printf("Unable to reset user table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	// drop records from db
	err := cfg.db.ResetUsersTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset user table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset users table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest CreateUserRequest
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

	// CreateUserParams struct
	createUserParams := database.CreateUserParams{
		CreatedBy:      createUserRequest.CreatedBy,
		UpdatedBy:      createUserRequest.UpdatedBy,
		Email:          createUserRequest.Email,
		HashedPassword: hashedPassword,
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

// logs in user and provides tokens
func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&userLoginRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// ensure login items arent empty
	if userLoginRequest.Email == "" || userLoginRequest.RawPassword == "" {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	userRecord, err := cfg.db.GetUserByEmailWithPassword(r.Context(), userLoginRequest.Email)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// hash & check password
	err = auth.CheckPasswordHash(userLoginRequest.RawPassword, userRecord.HashedPassword)
	if err != nil {
		respondWithError(err, http.StatusForbidden, w)
		return
	}
	// password checked, removing from memory
	userLoginRequest.RawPassword = ""

	// check email
	userLoginRequest.Email = strings.ToLower(userLoginRequest.Email)
	userRecord.Email = strings.ToLower(userRecord.Email)
	if userRecord.Email != userLoginRequest.Email {
		respondWithError(err, http.StatusForbidden, w)
		return
	}

	// user can log in, generate tokens
	log.Printf("Successfully logged in user: %q", userLoginRequest.Email)
}

// promotes user to admin
func (cfg *apiConfig) promoteUserToAdminHandler(w http.ResponseWriter, r *http.Request) {
	var adminStatusRequest AdminStatusRequest
	err := json.NewDecoder(r.Body).Decode(&adminStatusRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// validate that id is a users id
	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check that user is not already admin
	if userRecord.IsAdmin {
		respondWithError(fmt.Errorf("user is already admin"), http.StatusBadRequest, w)
		return
	}

	// make user admin
	err = cfg.db.PromoteUserToAdminByID(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// successful
	w.WriteHeader(http.StatusNoContent)
}

// demotes user from admin
func (cfg *apiConfig) demoteUserToAdminHandler(w http.ResponseWriter, r *http.Request) {
	var adminStatusRequest AdminStatusRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&adminStatusRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// validate that id is a users id
	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check that user is not demoted was never promoted
	if !userRecord.IsAdmin {
		respondWithError(fmt.Errorf("user is already not-admin"), http.StatusBadRequest, w)
		return
	}

	// demote user
	err = cfg.db.DemoteUserFromAdminByID(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// successful
	w.WriteHeader(http.StatusNoContent)
}
