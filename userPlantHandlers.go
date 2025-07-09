package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// === request / response structs ===

type UserCreatePlantRequest struct {
	PlantSpeciesID uuid.UUID  `json:"plantSpeciesID"`
	AdoptionDate   *time.Time `json:"adoptionDate"`
	Name           *string    `json:"plantName"`
}

type UserCreatePlantResponse struct {
	UsersPlantID     uuid.UUID  `json:"id"`
	PlantSpeciesID   uuid.UUID  `json:"plantSpeciesID"`
	PlantSpeciesName string     `json:"plantSpeciesName"`
	AdoptionDate     *time.Time `json:"adoptionDate,omitempty"`
	Name             *string    `json:"plantName,omitempty"`
}

type UserViewPlantResponse struct {
	UsersPlantID     uuid.UUID  `json:"id"`
	PlantSpeciesID   uuid.UUID  `json:"plantSpeciesID"`
	PlantSpeciesName string     `json:"plantSpeciesName"`
	AdoptionDate     *time.Time `json:"adoptionDate,omitempty"`
	Name             *string    `json:"plantName,omitempty"`
}

type UserUpdatePlantRequest struct {
	AdoptionDate *time.Time `json:"adoptionDate"`
	Name         *string    `json:"plantName"`
}

// performs authentication flow for normal users
func (cfg *apiConfig) userTokenAuthFlow(header http.Header, _ http.ResponseWriter) (uuid.UUID, error) {
	accessTokenProvided, err := auth.GetBearerToken(header, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get token from headers", "error", err)
		return uuid.Nil, nil
	}

	userID, err := auth.ValidateJWT(accessTokenProvided, cfg.JWTSecret, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		return uuid.Nil, nil
	}
	return userID, nil
}

// requires access token in auth header
// creates a user_plant
func (cfg *apiConfig) usersPlantsCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.userTokenAuthFlow(r.Header, w)
	if err != nil {
		cfg.sl.Warn("Could not complete user token auth flow during user request")
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// decode request body
	var createRequest UserCreatePlantRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode request body")
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// check request body
	if createRequest.PlantSpeciesID == uuid.Nil {
		cfg.sl.Debug("Request missing plant species id property")
		respondWithError(errors.New("no name property provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	// convert to database params
	adoptionDate := sql.NullTime{}
	plantName := sql.NullString{}

	if createRequest.AdoptionDate == nil {
		adoptionDate.Valid = false
	} else {
		adoptionDate.Valid = true
		adoptionDate.Time = *createRequest.AdoptionDate
	}
	if createRequest.Name == nil {
		plantName.Valid = false
	} else {
		plantName.Valid = true
		plantName.String = *createRequest.Name
	}

	createParams := database.CreateUsersPlantsParams{
		CreatedBy:    requestUserID,
		PlantID:      createRequest.PlantSpeciesID,
		UserID:       requestUserID,
		AdoptionDate: adoptionDate,
		Name:         plantName,
	}

	// perform database update
	userPlantRecord, err := cfg.db.CreateUsersPlants(r.Context(), createParams)
	if err != nil {
		cfg.sl.Debug("Could not create user's plant record", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// convert back to struct
	createResponse := UserCreatePlantResponse{
		UsersPlantID:     userPlantRecord.UsersPlantID,
		PlantSpeciesID:   userPlantRecord.SpeciesID,
		PlantSpeciesName: userPlantRecord.PlantSpeciesName,
		AdoptionDate:     createRequest.AdoptionDate,
		Name:             createRequest.Name,
	}

	// perform response
	cfg.sl.Debug("User successfully created a new users plant", "user id", requestUserID, "users plant id", createResponse.UsersPlantID, "plant species id", createRequest.PlantSpeciesID, "plant species name", userPlantRecord.PlantSpeciesName)
	respondWithJSON(http.StatusCreated, createResponse, w, cfg.sl)
}

// requires access token in auth header
// returns the users list of plants
func (cfg *apiConfig) usersPlantsListHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenProvided, err := auth.GetBearerToken(r.Header, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get token from headers", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	requestUserID, err := auth.ValidateJWT(accessTokenProvided, cfg.JWTSecret, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// get list of plants in user_plants table
	usersPlants, err := cfg.db.GetAllUsersPlantsOrderedByUpdated(r.Context(), requestUserID)
	if err != nil {
		cfg.sl.Debug("Could not get users' plant from database", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	if len(usersPlants) <= 0 {
		cfg.sl.Debug("User successfully requested their empty plants list", "user id", requestUserID)
		respondWithJSON(http.StatusOK, usersPlants, w, cfg.sl)
		return
	}

	// convert from database type to response type
	viewResponse := make([]UserViewPlantResponse, 0)
	for _, oldRecord := range usersPlants {
		var adoptionDate *time.Time
		var plantName *string

		if oldRecord.AdoptionDate.Valid {
			adoptionDate = &oldRecord.AdoptionDate.Time
		}
		if oldRecord.PlantName.Valid {
			plantName = &oldRecord.PlantName.String
		}

		newResponse := UserViewPlantResponse{
			UsersPlantID:     oldRecord.UsersPlantID,
			PlantSpeciesID:   oldRecord.PlantSpeciesID,
			PlantSpeciesName: oldRecord.SpeciesName,
			AdoptionDate:     adoptionDate,
			Name:             plantName,
		}
		viewResponse = append(viewResponse, newResponse)
	}

	cfg.sl.Debug("User successfully listed their users plants", "user id", requestUserID)
	respondWithJSON(http.StatusOK, viewResponse, w, cfg.sl)
}

func (cfg *apiConfig) userPlantsUpdateHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenProvided, err := auth.GetBearerToken(r.Header, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get token from headers", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	requestUserID, err := auth.ValidateJWT(accessTokenProvided, cfg.JWTSecret, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	plantIDStr := r.PathValue("plantID")
	plantID, err := uuid.Parse(plantIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant type id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var updateRequest UserUpdatePlantRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	if updateRequest.AdoptionDate == nil && updateRequest.Name == nil {
		cfg.sl.Debug("No updates provided in request")
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	newAdoptionDate := sql.NullTime{}
	newName := sql.NullString{}

	if updateRequest.AdoptionDate == nil {
		newAdoptionDate.Valid = false
	} else {
		newAdoptionDate.Valid = true
		newAdoptionDate.Time = *updateRequest.AdoptionDate
	}
	if updateRequest.Name == nil {
		newName.Valid = false
	} else {
		newName.Valid = true
		newName.String = *updateRequest.Name
	}

	updateParams := database.UpdateUsersPlantByIDParams{
		ID:           plantID,
		UpdatedBy:    requestUserID,
		AdoptionDate: newAdoptionDate,
		Name:         newName,
	}
	err = cfg.db.UpdateUsersPlantByID(r.Context(), updateParams)
	if err != nil {
		cfg.sl.Debug("Could not update users plant record", "error", err, "users plant id", plantID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("User successfully updated users plant", "user id", requestUserID, "users plant id", plantID)
	w.WriteHeader(http.StatusNoContent)
}
