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
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types ===

// create request for plant species
type AdminPlantSpeciesCreateRequest struct {
	SpeciesName      string `json:"speciesName"`
	HumanPoisonToxic *bool  `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic   *bool  `json:"petPoisonToxic,omitempty"`
	HumanEdible      *bool  `json:"humanEdible,omitempty"`
	PetEdible        *bool  `json:"petEdible,omitempty"`
}

// only provides client and updatable information
type AdminPlantSpeciesUpdateRequest struct {
	HumanPoisonToxic *bool `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic   *bool `json:"petPoisonToxic,omitempty"`
	HumanEdible      *bool `json:"humanEdible,omitempty"`
	PetEdible        *bool `json:"petEdible,omitempty"`
}

type AdminPlantSpeciesViewResponse struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	CreatedBy        uuid.UUID `json:"createdBy"`
	UpdatedBy        uuid.UUID `json:"updatedBy"`
	SpeciesName      string    `json:"speciesName"`
	HumanPoisonToxic *bool     `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic   *bool     `json:"petPoisonToxic,omitempty"`
	HumanEdible      *bool     `json:"humanEdible,omitempty"`
	PetEdible        *bool     `json:"petEdible,omitempty"`
}

// === handler functions ===

// GET json
func (cfg *apiConfig) adminPlantSpeciesViewHandler(w http.ResponseWriter, r *http.Request) {
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

	if len(plantSpeciesRecords) <= 0 {
		log.Print("Admin listed empty plant species list successfully.")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// NOTE: not the most effecient way to convert
	plantSpeciesResponse := make([]AdminPlantSpeciesViewResponse, 0)
	for _, oldRecord := range plantSpeciesRecords {
		// converting sql.NullBool to bool reference
		var humanPT *bool
		var humanE *bool
		var petPT *bool
		var petE *bool

		if oldRecord.HumanPoisonToxic.Valid {
			humanPT = &oldRecord.HumanPoisonToxic.Bool
		} else {
			humanPT = nil
		}
		if oldRecord.HumanEdible.Valid {
			humanE = &oldRecord.HumanEdible.Bool
		} else {
			humanE = nil
		}
		if oldRecord.PetPoisonToxic.Valid {
			petPT = &oldRecord.PetPoisonToxic.Bool
		} else {
			petPT = nil
		}
		if oldRecord.PetEdible.Valid {
			petE = &oldRecord.PetEdible.Bool
		} else {
			petE = nil
		}

		newResponse := AdminPlantSpeciesViewResponse{
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

// POST /api/v1/admin/plants/{plant_species_id}
func (cfg *apiConfig) adminReplacePlantSpeciesInfoHandler(w http.ResponseWriter, r *http.Request) {
	plantSpeciesIDStr := r.PathValue("plant_species_id")
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		log.Printf("Could not parse plant species id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check header for admin access token
	requestUserID, err := cfg.authorizeNormalAdmin(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var updateRequest AdminPlantSpeciesUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// converting from bool reference to sql.NullBool
	humanPT := sql.NullBool{}
	petPT := sql.NullBool{}
	humanE := sql.NullBool{}
	petE := sql.NullBool{}

	if updateRequest.HumanPoisonToxic == nil {
		humanPT.Valid = false
	} else {
		humanPT.Valid = true
		humanPT.Bool = *updateRequest.HumanPoisonToxic
	}
	if updateRequest.PetPoisonToxic == nil {
		petPT.Valid = false
	} else {
		petPT.Valid = true
		petPT.Bool = *updateRequest.PetPoisonToxic
	}
	if updateRequest.HumanEdible == nil {
		humanE.Valid = false
	} else {
		humanE.Valid = true
		humanE.Bool = *updateRequest.HumanEdible
	}
	if updateRequest.PetEdible == nil {
		petE.Valid = false
	} else {
		petE.Valid = true
		petE.Bool = *updateRequest.PetEdible
	}

	updateRequestParams := database.UpdatePlantSpeciesPropertiesByIDParams{
		ID:               plantSpeciesID,
		UpdatedBy:        requestUserID,
		HumanPoisonToxic: humanPT,
		PetPoisonToxic:   petPT,
		HumanEdible:      humanE,
		PetEdible:        petE,
	}

	err = cfg.db.UpdatePlantSpeciesPropertiesByID(r.Context(), updateRequestParams)
	if err != nil {
		log.Printf("Could not update plant species record %q due to: %q", plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q updated plant species %q successfully.", requestUserID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/admin/plants/{plant_species_id}
func (cfg *apiConfig) adminDeletePlantSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	plantSpeciesIDStr := r.PathValue("plant_species_id")
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		log.Printf("Could not parse plant species id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check header for admin access token
	requestUserID, err := cfg.authorizeNormalAdmin(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	requestUserNullUUID := uuid.NullUUID{
		UUID:  requestUserID,
		Valid: true,
	}

	deleteRequestParams := database.MarkPlantSpeciesAsDeletedByIDParams{
		ID:        plantSpeciesID,
		DeletedBy: requestUserNullUUID,
	}

	err = cfg.db.MarkPlantSpeciesAsDeletedByID(r.Context(), deleteRequestParams)
	if err != nil {
		log.Printf("Could not mark plant species %q as deleted due to: %q", plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully marked plant species %q as deleted.", requestUserID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}

// POST json to create plant
func (cfg *apiConfig) adminPlantSpeciesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.authorizeNormalAdmin(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var createRequest AdminPlantSpeciesCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// check all of the request properties
	if createRequest.SpeciesName == "" {
		log.Print("Request body was missing species name.")
		respondWithError(errors.New("no species name provided"), http.StatusBadRequest, w)
		return
	}

	// converting from bool reference to sql.NullBool
	humanPT := sql.NullBool{}
	petPT := sql.NullBool{}
	humanE := sql.NullBool{}
	petE := sql.NullBool{}

	if createRequest.HumanPoisonToxic == nil {
		// log.Print("Warning: human_poison_toxic field missing.")
		humanPT.Valid = false
	} else {
		// log.Print("Request body has human_poison_toxic is present.")
		humanPT.Valid = true
		humanPT.Bool = *createRequest.HumanPoisonToxic
	}
	if createRequest.PetPoisonToxic == nil {
		// log.Print("Warning: pet_poison_toxic field missing.")
		petPT.Valid = false
	} else {
		// log.Print("Request body has pet_poison_toxic is present.")
		petPT.Valid = true
		petPT.Bool = *createRequest.PetPoisonToxic
	}
	if createRequest.HumanEdible == nil {
		// log.Print("Warning: human_edible field missing.")
		humanE.Valid = false
	} else {
		// log.Print("Request body has human_edible is present.")
		humanE.Valid = true
		humanE.Bool = *createRequest.HumanEdible
	}
	if createRequest.PetEdible == nil {
		// log.Print("Warning: pet_edible field missing.")
		petE.Valid = false
	} else {
		// log.Print("Request body has pet_edible is present.")
		petE.Valid = true
		petE.Bool = *createRequest.PetEdible
	}

	createRequestParams := database.CreatePlantSpeciesParams{
		CreatedBy:        requestUserID,
		UpdatedBy:        requestUserID,
		SpeciesName:      createRequest.SpeciesName,
		HumanPoisonToxic: humanPT,
		PetPoisonToxic:   petPT,
		HumanEdible:      humanE,
		PetEdible:        petE,
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
