package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types

type AdminPlantsCreateRequest struct {
	Client           string `json:"client"`
	SpeciesName      string `json:"speciesName"`
	HumanPoisonToxic *bool  `json:"humanPoisonToxic"`
	PetPoisonToxic   *bool  `json:"petPoisonToxic"`
	HumanEdible      *bool  `json:"humanEdible"`
	PetEdible        *bool  `json:"petEdible"`
}

// === helper functions ===

// check header for admin access token
func (cfg *apiConfig) authorizeAdmin(r *http.Request) (uuid.UUID, error) {
	requestAccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		return uuid.UUID{}, err
	}

	requestUserID, err := auth.ValidateJWT(requestAccessToken, cfg.JWTSecret)
	if err != nil {
		return uuid.UUID{}, err
	}

	return requestUserID, nil
}

// === handler functions ===

// POST no request?
func (cfg *apiConfig) resetPlantSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	// ensure development platform
	if cfg.platform == "production" || cfg.platform == "" {
		log.Printf("Unable to reset user table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	// drop records from db
	err := cfg.db.ResetPlantSpeciesTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset plant_species table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset plant_species table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

// GET json
func (cfg *apiConfig) adminPlantsViewHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.authorizeAdmin(r)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), requestUserID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	if !userRecord.IsAdmin {
		respondWithError(fmt.Errorf("unauthorized"), http.StatusUnauthorized, w)
		return
	}
	// user is now authenticated below here

	plantSpeciesRecords, err := cfg.db.GetAllPlantSpeciesOrderedByCreated(r.Context())
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// log.Printf("plants list: %#v", plantSpeciesRecords)

	if len(plantSpeciesRecords) <= 0 {
		log.Printf("Admin listed empty plant species list successfully.")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	plantSpeciesData, err := json.Marshal(plantSpeciesRecords)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin listed plant species list successfully.")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(plantSpeciesData)
}

// POST json to create plant
func (cfg *apiConfig) adminAllInfoPlantsCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.authorizeAdmin(r)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var createRequest AdminPlantsCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// check all of the request properties
	if createRequest.Client == "" {
		respondWithError(errors.New("no client name provided"), http.StatusBadRequest, w)
		return
	}
	if createRequest.SpeciesName == "" {
		respondWithError(errors.New("no species name provided"), http.StatusBadRequest, w)
		return
	}

	// NOTE: convert to a warning if a single field is missing?
	humanPoisonToxic := sql.NullBool{}
	if createRequest.HumanPoisonToxic == nil {
		log.Print("WARNING: human_poison_toxic field missing.")
		humanPoisonToxic.Valid = false
	} else {
		humanPoisonToxic.Valid = true
		humanPoisonToxic.Bool = *createRequest.HumanPoisonToxic
	}

	petPoisonToxic := sql.NullBool{}
	if createRequest.PetPoisonToxic == nil {
		log.Print("WARNING: pet_poison_toxic field missing.")
		petPoisonToxic.Valid = false
	} else {
		petPoisonToxic.Valid = true
		petPoisonToxic.Bool = *createRequest.PetPoisonToxic
	}

	humanEdible := sql.NullBool{}
	if createRequest.HumanEdible == nil {
		log.Print("WARNING: human_edible field missing.")
		humanEdible.Valid = false
	} else {
		humanEdible.Valid = true
		humanEdible.Bool = *createRequest.HumanEdible
	}

	petEdible := sql.NullBool{}
	if createRequest.PetEdible == nil {
		log.Print("WARNING: pet_edible field missing.")
		petEdible.Valid = false
	} else {
		petEdible.Valid = true
		petEdible.Bool = *createRequest.PetEdible
	}

	createRequestParams := database.CreatePlantSpeciesParams{
		CreatedBy:        createRequest.Client,
		UpdatedBy:        createRequest.Client,
		SpeciesName:      createRequest.SpeciesName,
		HumanPoisonToxic: humanPoisonToxic,
		PetPoisonToxic:   petPoisonToxic,
		HumanEdible:      humanEdible,
		PetEdible:        petEdible,
	}
	_, err = cfg.db.CreatePlantSpecies(r.Context(), createRequestParams)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("User %q created a plant species %q", requestUserID, createRequest.SpeciesName)
	w.WriteHeader(http.StatusNoContent)
}
