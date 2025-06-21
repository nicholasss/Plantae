package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types

// promote / demote request
type AdminStatusRequest struct {
	ID uuid.UUID `json:"id"`
}

// register endpoint
type CreateUserRequest struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}

// login endpoint
type UserLoginRequest struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}
type UserLoginResponse struct {
	ID                    uuid.UUID `json:"id"`
	IsAdmin               bool      `json:"isAdmin"`
	AccessToken           string    `json:"token"`
	AccessTokenExpiresAt  time.Time `json:"tokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

// refresh endpoint
type AuthRefreshResponse struct {
	ID                   uuid.UUID `json:"id"`
	AccessToken          string    `json:"token"`
	AccessTokenExpiresAt time.Time `json:"tokenExpiresAt"`
}

type AuthRevokeRequest struct {
	ID uuid.UUID `json:"id"`
}

// === User Handler Functions ===

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest CreateUserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createUserRequest)
	if err != nil {
		// respond with error
	}
	// log.Print("Decoded create user request...")

	// check request params
	if createUserRequest.Email == "" {
		log.Print("Email received was empty.")
		respondWithError(nil, http.StatusBadRequest, w)
		return
	}
	if createUserRequest.RawPassword == "" {
		log.Print("Password received was empty.")
		respondWithError(nil, http.StatusBadRequest, w)
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(createUserRequest.RawPassword)
	createUserRequest.RawPassword = "" // GC collection
	if err != nil {
		log.Printf("Error hashing password for creating a user: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// user uuid generation
	newUserUUID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("Unable to create a UUID for a user due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// CreateUserParams struct
	createUserParams := database.CreateUserParams{
		ID:             newUserUUID,
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

	log.Printf("User %q was registered successfully.", userRecord.ID)
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
		log.Print("Login attempt failed with mis-matching password hashes.")
		respondWithError(err, http.StatusForbidden, w)
		return
	}
	// password checked, removing from memory
	userLoginRequest.RawPassword = ""

	// check email
	userLoginRequest.Email = strings.ToLower(userLoginRequest.Email)
	userRecord.Email = strings.ToLower(userRecord.Email)
	if userRecord.Email != userLoginRequest.Email {
		log.Printf("Login attempt failed for %q with email %q", userRecord.Email, userLoginRequest.Email)
		respondWithError(err, http.StatusForbidden, w)
		return
	}

	// user logged in, generate tokens
	// log.Printf("Successfully logged in user: %q", userRecord.ID)
	log.Printf("Generating tokens for user...")

	// refresh token
	userRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// store userRefreshToken in database
	refreshTokenExpiresAt := time.Now().Add(cfg.refreshTokenDuration)
	createRefreshToken := database.CreateRefreshTokenParams{
		RefreshToken: userRefreshToken,
		CreatedBy:    userRecord.ID,
		ExpiresAt:    refreshTokenExpiresAt,
	}

	// TODO: check for prexisting token, if exists then revoke it and replace
	_, err = cfg.db.CreateRefreshToken(r.Context(), createRefreshToken)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// access token
	accessTokenExpiresAt := time.Now().Add(cfg.accessTokenDuration)
	userAccessToken, err := auth.MakeJWT(userRecord.ID, cfg.JWTSecret, cfg.accessTokenDuration)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// json response
	userLoginResponse := UserLoginResponse{
		ID:                    userRecord.ID,
		IsAdmin:               userRecord.IsAdmin,
		AccessToken:           userAccessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          userRefreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}
	userLoginResponseData, err := json.Marshal(userLoginResponse)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	if platformNotProduction(cfg) {
		// log.Printf("DEBUG: User logged in: %q, accessToken: %q, refreshToken: %q", userLoginResponse.ID, userLoginResponse.AccessToken, userLoginResponse.RefreshToken)
	}

	log.Printf("User %q successfully logged in.", userRecord.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userLoginResponseData)
}

// accepts refresh token as authentication
// responds with a new access token if authorized
func (cfg *apiConfig) refreshUserHandler(w http.ResponseWriter, r *http.Request) {
	providedRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	refreshTokenRecord, err := cfg.db.GetUserFromRefreshToken(r.Context(), providedRefreshToken)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	if refreshTokenRecord.RevokedAt.Valid {
		// is marked as revoked
		if time.Now().After(refreshTokenRecord.RevokedAt.Time) {
			log.Print("Refresh token sent to POST /api/v1/auth/refresh was revoked.")
			respondWithError(err, http.StatusUnauthorized, w)
			return
		}

		// has been marked as revoked but in the future
		// these tokens should not be accepted
		// this may present a bug
		log.Print("!!! potential bug, check POST /api/refresh handler")
		log.Print("Refresh token will be revoked in the future.")
		respondWithError(errors.New("potential error in refresh token database. token revoked in the future"), http.StatusInternalServerError, w)
		return
	}

	if time.Now().UTC().After(refreshTokenRecord.ExpiresAt) {
		respondWithError(errors.New("bad request"), http.StatusBadRequest, w)
		return
	}

	accessTokenExpiresAt := time.Now().UTC().Add(cfg.accessTokenDuration)
	newAccessToken, err := auth.MakeJWT(refreshTokenRecord.UserID, cfg.JWTSecret, cfg.accessTokenDuration)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	refreshResponse := &AuthRefreshResponse{
		ID:                   refreshTokenRecord.UserID,
		AccessToken:          newAccessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}
	refreshResponseData, err := json.Marshal(refreshResponse)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("User %q	refreshed their access token successfully.", refreshTokenRecord.UserID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(refreshResponseData)
}

// accepts refresh token as authentication
// responds with 204 No Content if successfully revoked
func (cfg *apiConfig) revokeUserHandler(w http.ResponseWriter, r *http.Request) {
	providedRefreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var revokeRequest AuthRevokeRequest
	err = json.NewDecoder(r.Body).Decode(&revokeRequest)
	defer r.Body.Close()
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	revokeRefreshTokenParams := database.RevokeRefreshTokenWithTokenParams{
		RefreshToken: providedRefreshToken,
		UpdatedBy:    revokeRequest.ID,
	}
	revokeRecordUserID, err := cfg.db.RevokeRefreshTokenWithToken(r.Context(), revokeRefreshTokenParams)
	if err != nil {
		respondWithError(err, http.StatusUnauthorized, w)
		return
	}

	// token was revoked
	log.Printf("User %q revoked their refresh token successfully.", revokeRecordUserID)
	w.WriteHeader(http.StatusNoContent)
}
