package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/database"
)

// There should only be four types of light:
// - Bright direct
// - Bright indirect
// - Medium indirect
// - Low indirect

// === request response types ===

type AdminLightCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AdminLightUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// === handler functions ===

// POST /admin/light
func (cfg *apiConfig) adminLightCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var createRequest AdminLightCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// checking body properties
	if createRequest.Name == "" {
		log.Print("Request Body missing name property.")
		respondWithError(errors.New("no name provided"), http.StatusBadRequest, w)
		return
	}
	if createRequest.Description == "" {
		log.Print("Request body missing description property.")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w)
		return
	}

	createParams := database.CreateLightNeedParams{
		CreatedBy:   requestUserID,
		Name:        createRequest.Name,
		Description: createRequest.Description,
	}
	lightRecord, err := cfg.db.CreateLightNeed(r.Context(), createParams)
	if err != nil {
		log.Printf("Could not create light needs record due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully created lighting %q", requestUserID, lightRecord.ID)
	w.WriteHeader(http.StatusNoContent)
}

// GET /admin/light
func (cfg *apiConfig) adminLightViewHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	lightRecords, err := cfg.db.GetAllLightNeedsOrderedByCreated(r.Context())
	if err != nil {
		log.Printf("Could not get light needs records due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	lightData, err := json.Marshal(lightRecords)
	if err != nil {
		log.Printf("Could not marshal records to json due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q listed light needs successfully.", requestUserID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(lightData)
}

// PUT /admin/light/{lightID}
func (cfg *apiConfig) adminLightUpdateHandler(w http.ResponseWriter, r *http.Request) {
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		log.Printf("Could not parse light id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var updateRequest AdminLightUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// checking body properties
	if updateRequest.Name == "" {
		log.Print("Request Body missing name property.")
		respondWithError(errors.New("no name provided"), http.StatusBadRequest, w)
		return
	}
	if updateRequest.Description == "" {
		log.Print("Request body missing description property.")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w)
		return
	}

	updateParams := database.UpdateLightNeedsByIDParams{
		ID:          lightID,
		UpdatedBy:   requestUserID,
		Name:        updateRequest.Name,
		Description: updateRequest.Description,
	}
	err = cfg.db.UpdateLightNeedsByID(r.Context(), updateParams)
	if err != nil {
		log.Printf("Could not update light needs record due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully updated light needs record %q", requestUserID, lightID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /admin/light/{lightID}
func (cfg *apiConfig) adminLightDeleteHandler(w http.ResponseWriter, r *http.Request) {
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		log.Printf("Could not parse light id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	nullUserID := uuid.NullUUID{Valid: true, UUID: requestUserID}
	deleteParams := database.MarkLightNeedAsDeletedByIDParams{
		ID:        lightID,
		DeletedBy: nullUserID,
	}
	err = cfg.db.MarkLightNeedAsDeletedByID(r.Context(), deleteParams)
	if err != nil {
		log.Printf("Could not delete light needs record due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully deleted light needs record %q", requestUserID, lightID)
	w.WriteHeader(http.StatusNoContent)
}

// POST /admin/light/link{light id} ? plant species id = uuid
func (cfg *apiConfig) adminSetPlantAsLightNeedHandler(w http.ResponseWriter, r *http.Request) {
	// light id
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		log.Printf("Could not parse light id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		log.Print("No plant species id was specified in url query")
		respondWithError(errors.New("no plant species id was provided"), http.StatusBadRequest, w)
		return
	}
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		log.Printf("Could not parse plant species id from url query due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// user id
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// perform set
	lightNullID := uuid.NullUUID{
		UUID:  lightID,
		Valid: true,
	}
	setParams := database.SetPlantSpeciesAsLightNeedParams{
		ID:           plantSpeciesID,
		LightNeedsID: lightNullID,
		UpdatedBy:    requestUserID,
	}
	err = cfg.db.SetPlantSpeciesAsLightNeed(r.Context(), setParams)
	if err != nil {
		log.Printf("Could not set light need %q for plant species %q due to: %q", lightID, plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully set light need %q for species %q ", requestUserID, lightID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /admin/light/link{light id} ? plant species id = uuid
func (cfg *apiConfig) adminUnsetPlantAsLightNeedHandler(w http.ResponseWriter, r *http.Request) {
	// light id
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		log.Printf("Could not parse light id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		log.Print("No plant species id was specified in url query")
		respondWithError(errors.New("no plant species id was provided"), http.StatusBadRequest, w)
		return
	}
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		log.Printf("Could not parse plant species id from url query due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// user id
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// perform unset
	unsetParams := database.UnsetPlantSpeciesAsLightNeedParams{
		ID:        plantSpeciesID,
		UpdatedBy: requestUserID,
	}
	err = cfg.db.UnsetPlantSpeciesAsLightNeed(r.Context(), unsetParams)
	if err != nil {
		log.Printf("Could not unset light need %q for plant species %q due to: %q", lightID, plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully unset light need %q for species %q ", requestUserID, lightID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}
