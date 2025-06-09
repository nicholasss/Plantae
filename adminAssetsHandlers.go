package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types

type AdminPlantsCreateRequest struct {
	Client           string `json:"client"`
	SpeciesName      string `json:"speciesName"`
	HumanPoisonToxic *bool  `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic   *bool  `json:"petPoisonToxic,omitempty"`
	HumanEdible      *bool  `json:"humanEdible,omitempty"`
	PetEdible        *bool  `json:"petEdible,omitempty"`
}

type AdminPlantsViewResponse struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedBy        string    `json:"createdBy"`
	UpdatedBy        string    `json:"updatedBy"`
	SpeciesName      string    `json:"speciesName"`
	HumanPoisonToxic *bool     `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic   *bool     `json:"petPoisonToxic,omitempty"`
	HumanEdible      *bool     `json:"humanEdible,omitempty"`
	PetEdible        *bool     `json:"petEdible,omitempty"`
}

// === helper functions ===

// check header for admin access token
func (cfg *apiConfig) authorizeNormalAdmin(r *http.Request) (uuid.UUID, error) {
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
	if platformProduction(cfg) {
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
	requestUserID, err := cfg.authorizeNormalAdmin(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), requestUserID)
	if err != nil {
		log.Printf("Could not get user record via user id due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	if !userRecord.IsAdmin {
		log.Print("Could not view plants due to user not being admin.")
		respondWithError(fmt.Errorf("unauthorized"), http.StatusUnauthorized, w)
		return
	}
	// user is now authenticated below here

	plantSpeciesRecords, err := cfg.db.GetAllPlantSpeciesOrderedByCreated(r.Context())
	if err != nil {
		log.Printf("Could not get plant species records due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// log.Printf("plants list: %#v", plantSpeciesRecords)

	if len(plantSpeciesRecords) <= 0 {
		log.Print("Admin listed empty plant species list successfully.")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	plantSpeciesResponse := make([]AdminPlantsViewResponse, 0)
	for _, oldRecord := range plantSpeciesRecords {
		var humanPT *bool
		if oldRecord.HumanPoisonToxic.Valid {
			humanPT = &oldRecord.HumanPoisonToxic.Bool
		} else {
			humanPT = nil
		}

		var humanE *bool
		if oldRecord.HumanEdible.Valid {
			humanE = &oldRecord.HumanEdible.Bool
		} else {
			humanE = nil
		}

		var petPT *bool
		if oldRecord.PetPoisonToxic.Valid {
			petPT = &oldRecord.PetPoisonToxic.Bool
		} else {
			petPT = nil
		}

		var petE *bool
		if oldRecord.PetEdible.Valid {
			petE = &oldRecord.PetEdible.Bool
		} else {
			petE = nil
		}

		newResponse := AdminPlantsViewResponse{
			ID:               oldRecord.ID,
			CreatedAt:        oldRecord.CreatedAt,
			UpdatedAt:        oldRecord.UpdatedAt,
			CreatedBy:        oldRecord.CreatedBy,
			UpdatedBy:        oldRecord.UpdatedBy,
			SpeciesName:      oldRecord.SpeciesName,
			HumanPoisonToxic: humanPT,
			HumanEdible:      humanE,
			PetPoisonToxic:   petPT,
			PetEdible:        petE,
		}

		plantSpeciesResponse = append(plantSpeciesResponse, newResponse)
	}

	plantSpeciesData, err := json.Marshal(plantSpeciesResponse)
	if err != nil {
		log.Printf("Could not get plant_species records due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin listed plant species list successfully.")
	log.Printf("DEBUG: list of plants: %s", string(plantSpeciesData))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(plantSpeciesData)
}

// POST json to create plant
func (cfg *apiConfig) adminAllInfoPlantsCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.authorizeNormalAdmin(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var createRequest AdminPlantsCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// check all of the request properties
	if createRequest.Client == "" {
		log.Print("Request body was missing client field.")
		respondWithError(errors.New("no client name provided"), http.StatusBadRequest, w)
		return
	}
	if createRequest.SpeciesName == "" {
		log.Print("Request body was missing species name.")
		respondWithError(errors.New("no species name provided"), http.StatusBadRequest, w)
		return
	}

	// NOTE: convert to a warning if a single field is missing instead of individual?
	humanPoisonToxic := sql.NullBool{}
	if createRequest.HumanPoisonToxic == nil {
		log.Print("WARNING: human_poison_toxic field missing.")
		humanPoisonToxic.Valid = false
	} else {
		log.Print("Request body has human_poison_toxic is present.")
		humanPoisonToxic.Valid = true
		humanPoisonToxic.Bool = *createRequest.HumanPoisonToxic
	}

	petPoisonToxic := sql.NullBool{}
	if createRequest.PetPoisonToxic == nil {
		log.Print("WARNING: pet_poison_toxic field missing.")
		petPoisonToxic.Valid = false
	} else {
		log.Print("Request body has pet_poison_toxic is present.")
		petPoisonToxic.Valid = true
		petPoisonToxic.Bool = *createRequest.PetPoisonToxic
	}

	humanEdible := sql.NullBool{}
	if createRequest.HumanEdible == nil {
		log.Print("WARNING: human_edible field missing.")
		humanEdible.Valid = false
	} else {
		log.Print("Request body has human_edible is present.")
		humanEdible.Valid = true
		humanEdible.Bool = *createRequest.HumanEdible
	}

	petEdible := sql.NullBool{}
	if createRequest.PetEdible == nil {
		log.Print("WARNING: pet_edible field missing.")
		petEdible.Valid = false
	} else {
		log.Print("Request body has pet_edible is present.")
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
		log.Printf("Could not create plant_species record in database due to %q.", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q created a plant species %q.", requestUserID, createRequest.SpeciesName)
	w.WriteHeader(http.StatusNoContent)
}
