package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types ===

// AdminPlantSpeciesCreateRequest is for decoding plant species create requests.
type AdminPlantSpeciesCreateRequest struct {
	SpeciesName      string `json:"speciesName"`
	HumanPoisonToxic *bool  `json:"humanPoisonToxic"`
	PetPoisonToxic   *bool  `json:"petPoisonToxic"`
	HumanEdible      *bool  `json:"humanEdible"`
	PetEdible        *bool  `json:"petEdible"`
}

// AdminPlantSpeciesUpdateRequest is for decoding plant species update requests.
type AdminPlantSpeciesUpdateRequest struct {
	HumanPoisonToxic *bool `json:"humanPoisonToxic"`
	PetPoisonToxic   *bool `json:"petPoisonToxic"`
	HumanEdible      *bool `json:"humanEdible"`
	PetEdible        *bool `json:"petEdible"`
}

// AdminPlantSpeciesViewResponse is for responding to plant species view requests.
type AdminPlantSpeciesViewResponse struct {
	ID               uuid.UUID `json:"id"`
	SpeciesName      string    `json:"speciesName"`
	HumanPoisonToxic *bool     `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic   *bool     `json:"petPoisonToxic,omitempty"`
	HumanEdible      *bool     `json:"humanEdible,omitempty"`
	PetEdible        *bool     `json:"petEdible,omitempty"`
}

// === handler functions ===

// GET json
func (cfg *apiConfig) adminPlantSpeciesViewHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	plantSpeciesRecords, err := cfg.db.GetAllPlantSpeciesOrderedByCreated(r.Context())
	if err != nil {
		cfg.sl.Debug("Could not get plant species records", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	if len(plantSpeciesRecords) <= 0 {
		cfg.sl.Debug("Admin successfully listed empty plant species list", "admin id", requestUserID)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// TODO: not the most efficient way to convert, is there another way?
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
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully listed plant species list", "admin id", requestUserID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(plantSpeciesData)
}

// PUT /api/v1/admin/plants/{plantSpeciesID}
func (cfg *apiConfig) adminReplacePlantSpeciesInfoHandler(w http.ResponseWriter, r *http.Request) {
	plantSpeciesIDStr := r.PathValue("plantSpeciesID")
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse species id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var updateRequest AdminPlantSpeciesUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
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
		cfg.sl.Debug("Could not update plant species record", "error", err, "plant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully updated plant species", "admin id", requestUserID, "plant species id", plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/admin/plants/{plantSpeciesID}
func (cfg *apiConfig) adminDeletePlantSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	plantSpeciesIDStr := r.PathValue("plantSpeciesID")
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse species id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
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
		cfg.sl.Debug("Could not mark plant species as deleted", "error", err, "plant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully marked plant species as deleted", "admin id", requestUserID, "plant species id", plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}

// POST json to create plant
func (cfg *apiConfig) adminPlantSpeciesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var createRequest AdminPlantSpeciesCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// check all of the request properties
	if createRequest.SpeciesName == "" {
		cfg.sl.Debug("Request body missing species name")
		respondWithError(errors.New("no species name provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	// converting from bool reference to sql.NullBool
	humanPT := sql.NullBool{}
	petPT := sql.NullBool{}
	humanE := sql.NullBool{}
	petE := sql.NullBool{}

	if createRequest.HumanPoisonToxic == nil {
		humanPT.Valid = false
	} else {
		humanPT.Valid = true
		humanPT.Bool = *createRequest.HumanPoisonToxic
	}
	if createRequest.PetPoisonToxic == nil {
		petPT.Valid = false
	} else {
		petPT.Valid = true
		petPT.Bool = *createRequest.PetPoisonToxic
	}
	if createRequest.HumanEdible == nil {
		humanE.Valid = false
	} else {
		humanE.Valid = true
		humanE.Bool = *createRequest.HumanEdible
	}
	if createRequest.PetEdible == nil {
		petE.Valid = false
	} else {
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
	speciesRecord, err := cfg.db.CreatePlantSpecies(r.Context(), createRequestParams)
	if err != nil {
		cfg.sl.Debug("Could not create plant species in database", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully created plant species", "admin id", requestUserID, "species id", speciesRecord.ID, "species name", speciesRecord.SpeciesName)
	w.WriteHeader(http.StatusNoContent)
}
