package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types

// CreateUserRequest is for decoding create user requests.
type CreateUserRequest struct {
	Email        string `json:"email"`
	RawPassword  string `json:"password"`
	LangCodePref string `json:"langCodePref"`
}

// CreateUserResponse is for decoding create user requests.
type CreateUserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	LangCodePref string    `json:"langCodePref"`
	JoinDate     time.Time `json:"joinDate"`
	IsAdmin      bool      `json:"isAdmin"`
}

// UserLoginRequest is for decoding user login requests.
type UserLoginRequest struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}

// UserLoginResponse is for encoding user loging responses.
type UserLoginResponse struct {
	ID                    uuid.UUID `json:"id"`
	LangCodePref          string    `json:"langCodePref"`
	JoinDate              time.Time `json:"joinDate"`
	IsAdmin               bool      `json:"isAdmin"`
	AccessToken           string    `json:"token"`
	AccessTokenExpiresAt  time.Time `json:"tokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

// AuthRefreshResponse is for encoding user access token responses.
type AuthRefreshResponse struct {
	ID                   uuid.UUID `json:"id"`
	AccessToken          string    `json:"token"`
	AccessTokenExpiresAt time.Time `json:"tokenExpiresAt"`
}

// AuthRevokeRequest is for decoding user refresh token requests.
type AuthRevokeRequest struct {
	ID uuid.UUID `json:"id"`
}

// === User Handler Functions ===

// POST /api/v1/auth/register
func (cfg *apiConfig) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var createUserRequest CreateUserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&createUserRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// check request params
	if createUserRequest.Email == "" {
		cfg.sl.Debug("Request body missing email")
		respondWithError(nil, http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createUserRequest.RawPassword == "" {
		cfg.sl.Debug("Request body missing password")
		respondWithError(nil, http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createUserRequest.LangCodePref == "" {
		cfg.sl.Debug("Request body missing language preference")
		respondWithError(nil, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// language preference check
	userRequestedLangName, ok := LangCodes[createUserRequest.LangCodePref]
	if !ok {
		cfg.sl.Debug("Requested language code was not found", "lang code", createUserRequest.LangCodePref)
		respondWithError(errors.New("language code requested does not exist"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	cfg.sl.Debug("User is registering with language", "lang code", createUserRequest.LangCodePref, "lang name", userRequestedLangName)

	// hash password
	hashedPassword, err := auth.HashPassword(createUserRequest.RawPassword, cfg.sl)
	createUserRequest.RawPassword = "" // GC collection
	if err != nil {
		cfg.sl.Debug("Error hashing user's password", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// user uuid generation
	newUserUUID, err := uuid.NewUUID()
	if err != nil {
		cfg.sl.Debug("Unable to generate a new uuid for user", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// add user to database
	createUserParams := database.CreateUserParams{
		ID:             newUserUUID,
		LangCodePref:   createUserRequest.LangCodePref,
		Email:          createUserRequest.Email,
		HashedPassword: hashedPassword,
	}
	userRecord, err := cfg.db.CreateUser(r.Context(), createUserParams)
	if err != nil {
		cfg.sl.Debug("Could not create user in database", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	createUserResponse := CreateUserResponse{
		ID:           userRecord.ID,
		Email:        userRecord.Email,
		LangCodePref: userRecord.LangCodePref,
		JoinDate:     userRecord.JoinDate,
		IsAdmin:      userRecord.IsAdmin,
	}

	// return the userRecord without password
	cfg.sl.Debug("User successfully registered", "user id", userRecord.ID)
	respondWithJSON(http.StatusCreated, createUserResponse, w, cfg.sl)
}

// logs in user and provides tokens
// POST /api/v1/auth/login
// POST /login is an exception
// -- it typically responds with HTTP 200 and response
// TODO: check for prexisting token, if exists then revoke it and replace
func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	var userLoginRequest UserLoginRequest
	err := json.NewDecoder(r.Body).Decode(&userLoginRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// ensure login items arent empty
	if userLoginRequest.Email == "" || userLoginRequest.RawPassword == "" {
		cfg.sl.Debug("Request body missing email or password")
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	userRecord, err := cfg.db.GetUserByEmailWithPassword(r.Context(), userLoginRequest.Email)
	if err != nil {
		cfg.sl.Debug("Unable to retreive user record with email", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// hash & check password
	err = auth.CheckPasswordHash(userLoginRequest.RawPassword, userRecord.HashedPassword, cfg.sl)
	if err != nil {
		cfg.sl.Debug("User's login attempt failed due to mis-matching passwords", "error", err)
		respondWithError(err, http.StatusForbidden, w, cfg.sl)
		return
	}
	// password checked, removing from memory
	userLoginRequest.RawPassword = ""

	// check email
	userLoginRequest.Email = strings.ToLower(userLoginRequest.Email)
	userRecord.Email = strings.ToLower(userRecord.Email)
	if userRecord.Email != userLoginRequest.Email {
		cfg.sl.Debug("User's login attempt failed due to mis-matching email", "request email", userLoginRequest.Email, "account email", userRecord)
		cfg.sl.Warn("User record details may be mis-matching queried values")
		respondWithError(err, http.StatusForbidden, w, cfg.sl)
		return
	}

	// user logged in, generate tokens
	cfg.sl.Debug("User logged in, generating new tokens")

	// refresh token
	userRefreshToken, err := auth.MakeRefreshToken(cfg.sl)
	if err != nil {
		cfg.sl.Debug("Unable to create new refresh token for user", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// store userRefreshToken in database
	refreshTokenExpiresAt := time.Now().Add(cfg.refreshTokenDuration)
	createRefreshToken := database.CreateRefreshTokenParams{
		RefreshToken: userRefreshToken,
		CreatedBy:    userRecord.ID,
		ExpiresAt:    refreshTokenExpiresAt,
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), createRefreshToken)
	if err != nil {
		cfg.sl.Warn("Unable to put a user's new refresh token into database", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// access token
	accessTokenExpiresAt := time.Now().Add(cfg.accessTokenDuration)
	userAccessToken, err := auth.MakeJWT(userRecord.ID, cfg.JWTSecret, cfg.accessTokenDuration, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Unable to create a new access token for user's login", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// json response
	userLoginResponse := UserLoginResponse{
		ID:                    userRecord.ID,
		LangCodePref:          userRecord.LangCodePref,
		JoinDate:              userRecord.JoinDate,
		IsAdmin:               userRecord.IsAdmin,
		AccessToken:           userAccessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          userRefreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}

	if platformNotProduction(cfg) {
		cfg.sl.Debug("Listing user info at login", "user id", userLoginResponse.ID, "user access token", userLoginResponse.AccessToken, "user refresh token", userLoginResponse.RefreshToken)
	}

	cfg.sl.Debug("User successfully logged in", "user id", userRecord.ID)
	respondWithJSON(http.StatusOK, userLoginResponse, w, cfg.sl)
}

// accepts refresh token as authentication
// responds with a new access token if authorized
// POST /api/v1/auth/refresh
// POST /refresh is about changing state, not creating a new record on the server
func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	providedRefreshToken, err := auth.GetBearerToken(r.Header, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get refresh token from headers", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	refreshTokenRecord, err := cfg.db.GetUserFromRefreshToken(r.Context(), providedRefreshToken)
	if err != nil {
		cfg.sl.Debug("Could not get refresh token record from database", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	if refreshTokenRecord.RevokedAt.Valid {
		// is marked as revoked
		if time.Now().After(refreshTokenRecord.RevokedAt.Time) {
			cfg.sl.Debug("Refresh token was revoked for user", "user id", refreshTokenRecord.UserID)
			respondWithError(err, http.StatusUnauthorized, w, cfg.sl)
			return
		}

		// has been marked as revoked but in the future
		// these tokens should not be accepted
		// this may present a bug
		// one place to check is the POST /refresh handler
		cfg.sl.Warn("Refresh token is marked as revoked but is still valid.")
		respondWithError(errors.New("valid token will be revoked in the future"), http.StatusInternalServerError, w, cfg.sl)
		return
	}

	if time.Now().UTC().After(refreshTokenRecord.ExpiresAt) {
		cfg.sl.Debug("Refresh token has expired")
		respondWithError(errors.New("bad request"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	accessTokenExpiresAt := time.Now().UTC().Add(cfg.accessTokenDuration)
	newAccessToken, err := auth.MakeJWT(refreshTokenRecord.UserID, cfg.JWTSecret, cfg.accessTokenDuration, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not create a new access token", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	refreshResponse := &AuthRefreshResponse{
		ID:                   refreshTokenRecord.UserID,
		AccessToken:          newAccessToken,
		AccessTokenExpiresAt: accessTokenExpiresAt,
	}

	cfg.sl.Debug("User successfully refreshed their access token", "user id", refreshTokenRecord.UserID)
	respondWithJSON(http.StatusOK, refreshResponse, w, cfg.sl)
}

// accepts refresh token as authentication
// responds with 204 No Content if successfully revoked
// POST /api/v1/auth/revoke
// POST /revoke sending 204 No Content makes sense,
// -- as it is a clean response to changing an existing record on the server
func (cfg *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	providedRefreshToken, err := auth.GetBearerToken(r.Header, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get refresh token from headers", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var revokeRequest AuthRevokeRequest
	err = json.NewDecoder(r.Body).Decode(&revokeRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	revokeRefreshTokenParams := database.RevokeRefreshTokenWithTokenParams{
		RefreshToken: providedRefreshToken,
		UpdatedBy:    revokeRequest.ID,
	}
	revokeRecordUserID, err := cfg.db.RevokeRefreshTokenWithToken(r.Context(), revokeRefreshTokenParams)
	if err != nil {
		cfg.sl.Debug("Could not revoke refresh token in database", "error", err, "refresh token id", revokeRequest.ID)
		respondWithError(err, http.StatusUnauthorized, w, cfg.sl)
		return
	}

	// token was revoked
	cfg.sl.Debug("User successfully revoked their refresh token", "user	id", revokeRecordUserID)
	w.WriteHeader(http.StatusNoContent)
}
