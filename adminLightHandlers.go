package main

import (
	"encoding/json"
	"errors"
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

type AdminSetLightResponse struct {
	LightNeedID      uuid.UUID `json:"lightNeedID"`
	PlantSpeciesID   uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName string    `json:"plantSpeciesName"`
}

type AdminUnsetLightResponse struct {
	PlantSpeciesID   uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName string    `json:"plantSpeciesName"`
}

// === handler functions ===

// TODO: ensure resource is sent back
// POST /api/v1/admin/light
func (cfg *apiConfig) adminLightCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var createRequest AdminLightCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// checking body properties
	if createRequest.Name == "" {
		cfg.sl.Debug("Request midding name property")
		respondWithError(errors.New("no name property provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createRequest.Description == "" {
		cfg.sl.Debug("Request midding description property")
		respondWithError(errors.New("no description property provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	createParams := database.CreateLightNeedParams{
		CreatedBy:   requestUserID,
		Name:        createRequest.Name,
		Description: createRequest.Description,
	}
	_, err = cfg.db.CreateLightNeed(r.Context(), createParams)
	if err != nil {
		cfg.sl.Debug("Could not create record in light_needs", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully completed request", "admin id", requestUserID)
	w.WriteHeader(http.StatusNoContent)
}

// GET /admin/light
func (cfg *apiConfig) adminLightViewHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	lightRecords, err := cfg.db.GetAllLightNeedsOrderedByCreated(r.Context())
	if err != nil {
		cfg.sl.Debug("Could not get light_needs record", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	lightData, err := json.Marshal(lightRecords)
	if err != nil {
		cfg.sl.Debug("Could not marshal records to json", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully completed request", "admin id", requestUserID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(lightData)
}

// PUT /admin/light/{lightID}
func (cfg *apiConfig) adminLightUpdateHandler(w http.ResponseWriter, r *http.Request) {
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse lightID from url path", "error", err)
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

	var updateRequest AdminLightUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// checking body properties
	if updateRequest.Name == "" {
		cfg.sl.Debug("Request body missing name property")
		respondWithError(errors.New("no name provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	if updateRequest.Description == "" {
		cfg.sl.Debug("Request body missing description property")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w, cfg.sl)
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
		cfg.sl.Debug("Could not update light needs record", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully completed request", "admin id", requestUserID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /admin/light/{lightID}
func (cfg *apiConfig) adminLightDeleteHandler(w http.ResponseWriter, r *http.Request) {
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse light id from url path", "error", err)
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

	nullUserID := uuid.NullUUID{Valid: true, UUID: requestUserID}
	deleteParams := database.MarkLightNeedAsDeletedByIDParams{
		ID:        lightID,
		DeletedBy: nullUserID,
	}
	err = cfg.db.MarkLightNeedAsDeletedByID(r.Context(), deleteParams)
	if err != nil {
		cfg.sl.Debug("Could not delete light needs records", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully completed request", "admin id", requestUserID)
	w.WriteHeader(http.StatusNoContent)
}

// TODO: ensure resource is sent back
// POST /admin/light/link{light id} ? plant species id = uuid
func (cfg *apiConfig) adminSetPlantAsLightNeedHandler(w http.ResponseWriter, r *http.Request) {
	// light id
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse light id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		cfg.sl.Debug("No plant species id was specified in url query")
		respondWithError(errors.New("no plant species id was provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant species id from url query", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// user id
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
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
	lightRecord, err := cfg.db.SetPlantSpeciesAsLightNeed(r.Context(), setParams)
	if err != nil {
		cfg.sl.Debug("Could not set light need for plant species id", "error", err, "light need id", lightID, "plant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	lightResponse := AdminSetLightResponse{
		LightNeedID:      lightID,
		PlantSpeciesID:   plantSpeciesID,
		PlantSpeciesName: lightRecord.SpeciesName,
	}
	lightData, err := json.Marshal(&lightResponse)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully set plant species to light need", "admin id", requestUserID, "plant species id", plantSpeciesID, "light need", lightID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(lightData)
}

// DELETE /admin/light/link{light id} ? plant species id = uuid
func (cfg *apiConfig) adminUnsetPlantAsLightNeedHandler(w http.ResponseWriter, r *http.Request) {
	// light id
	lightIDStr := r.PathValue("lightID")
	lightID, err := uuid.Parse(lightIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse light id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		cfg.sl.Debug("No plant species id was specified in url query")
		respondWithError(errors.New("no plant species id was provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant species id from url query", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// user id
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// perform unset
	unsetParams := database.UnsetPlantSpeciesAsLightNeedParams{
		ID:        plantSpeciesID,
		UpdatedBy: requestUserID,
	}
	lightRecord, err := cfg.db.UnsetPlantSpeciesAsLightNeed(r.Context(), unsetParams)
	if err != nil {
		cfg.sl.Debug("Could not unset light need for plant species id", "error", err, "light need id", lightID, "plant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	lightResponse := AdminUnsetLightResponse{
		PlantSpeciesID:   plantSpeciesID,
		PlantSpeciesName: lightRecord.SpeciesName,
	}
	lightData, err := json.Marshal(&lightResponse)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully unset plant species to light need", "admin id", requestUserID, "plant species id", plantSpeciesID, "light need", lightID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(lightData)
}
