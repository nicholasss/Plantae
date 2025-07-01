package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/database"
)

// There should only be two different ways to water:
//
// x millimeters of soil should be dry before watering
// - Tropical
// - Temperate
//
// x number of days between watering
// - Semi-Arid
// - Arid

// === request response types ===

type AdminWaterCreateRequest struct {
	PlantType   string `json:"plantType"`
	Description string `json:"description"`
	DrySoilMM   *int32 `json:"drySoilMM"`
	DrySoilDays *int32 `json:"drySoilDays"`
}

type AdminWaterResponse struct {
	ID          uuid.UUID `json:"id"`
	PlantType   string    `json:"plantType"`
	Description string    `json:"description"`
	DrySoilMM   *int32    `json:"drySoilMM,omitempty"`
	DrySoilDays *int32    `json:"drySoilDays,omitempty"`
}

type AdminSetWaterResponse struct {
	WaterNeedID      uuid.UUID `json:"waterNeedID"`
	PlantSpeciesID   uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName string    `json:"plantSpeciesName"`
}

type AdminUnsetWaterResponse struct {
	PlantSpeciesID   uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName string    `json:"plantSpeciesName"`
}

// === handler functions ===

// TODO: ensure resource is sent back
// POST /api/v1/admin/water
func (cfg *apiConfig) adminWaterCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var createRequest AdminWaterCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// checking body
	if createRequest.PlantType == "" {
		cfg.sl.Debug("Request body missing plant type")
		respondWithError(errors.New("no plant type provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createRequest.Description == "" {
		cfg.sl.Debug("Request body missing description")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	mmRequest := strings.ToLower(createRequest.PlantType) == "tropical" || strings.ToLower(createRequest.PlantType) == "temperate"
	dayRequest := strings.ToLower(createRequest.PlantType) == "semi-arid" || strings.ToLower(createRequest.PlantType) == "arid"

	var waterResponse AdminWaterResponse

	if dayRequest {
		if createRequest.DrySoilDays == nil {
			cfg.sl.Debug("Request body missing dry soil days")
			respondWithError(errors.New("no dry soil days provided"), http.StatusBadRequest, w, cfg.sl)
			return
		}

		nullDrySoilDays := sql.NullInt32{Int32: *createRequest.DrySoilDays, Valid: true}
		createParams := database.CreateWaterDryDaysParams{
			CreatedBy:   requestUserID,
			PlantType:   createRequest.PlantType,
			Description: createRequest.Description,
			DrySoilDays: nullDrySoilDays,
		}
		waterRecord, err := cfg.db.CreateWaterDryDays(r.Context(), createParams)
		if err != nil {
			cfg.sl.Debug("Could not create water record", "error", err)
			respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
			return
		}

		waterResponse.ID = waterRecord.ID
		waterResponse.PlantType = waterRecord.PlantType
		waterResponse.Description = waterRecord.Description
		waterResponse.DrySoilDays = &waterRecord.DrySoilDays.Int32

	} else if mmRequest {
		if createRequest.DrySoilMM == nil {
			cfg.sl.Debug("Request body missing dry soil mm")
			respondWithError(errors.New("no dry soil mm provided"), http.StatusBadRequest, w, cfg.sl)
			return
		}

		nullDrySoilMM := sql.NullInt32{Int32: *createRequest.DrySoilMM, Valid: true}
		createParams := database.CreateWaterDryMMParams{
			CreatedBy:   requestUserID,
			PlantType:   createRequest.PlantType,
			Description: createRequest.Description,
			DrySoilMm:   nullDrySoilMM,
		}
		waterRecord, err := cfg.db.CreateWaterDryMM(r.Context(), createParams)
		if err != nil {
			cfg.sl.Debug("Could not create water record", "error", err)
			respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
			return
		}

		waterResponse.ID = waterRecord.ID
		waterResponse.PlantType = waterRecord.PlantType
		waterResponse.Description = waterRecord.Description
		waterResponse.DrySoilMM = &waterRecord.DrySoilMm.Int32

	} else {
		cfg.sl.Debug("Invalid plant type provided", "plant type", createRequest.PlantType)
		respondWithError(errors.New("invalid plant type"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	waterData, err := json.Marshal(&waterResponse)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully created water need", "admin id", requestUserID, "water need id", waterResponse.ID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(waterData)
}

func (cfg *apiConfig) adminWaterViewHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	waterRecords, err := cfg.db.GetAllWaterNeedsOrderedByCreated(r.Context())
	if err != nil {
		cfg.sl.Debug("Could not get water needs records", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// NOTE: may not be the most effificent
	var waterResponses []AdminWaterResponse
	for _, record := range waterRecords {
		mmRecord := strings.ToLower(record.PlantType) == "tropical" || strings.ToLower(record.PlantType) == "temperate"
		dayRecord := strings.ToLower(record.PlantType) == "semi-arid" || strings.ToLower(record.PlantType) == "arid"

		if mmRecord {
			drySoilMM := record.DrySoilMm.Int32
			newRecord := AdminWaterResponse{
				ID:          record.ID,
				PlantType:   record.PlantType,
				Description: record.Description,
				DrySoilMM:   &drySoilMM,
			}
			waterResponses = append(waterResponses, newRecord)
		} else if dayRecord {
			drySoilDays := record.DrySoilDays.Int32
			newRecord := AdminWaterResponse{
				ID:          record.ID,
				PlantType:   record.PlantType,
				Description: record.Description,
				DrySoilDays: &drySoilDays,
			}
			waterResponses = append(waterResponses, newRecord)
		} else {
			cfg.sl.Warn("Unknown plant type returned from database", "plant type id", record.ID, "plant type", record.PlantType)
		}
	}

	waterData, err := json.Marshal(waterResponses)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully listed water needs list", "admin id", requestUserID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(waterData)
}

// DELETE /admin/water/{waterID}
func (cfg *apiConfig) adminWaterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// water need record to mark as deleted
	waterIDStr := r.PathValue("waterID")
	waterID, err := uuid.Parse(waterIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse water id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// perform delete
	nullUserID := uuid.NullUUID{UUID: requestUserID, Valid: true}
	deleteParams := database.MarkWaterNeedAsDeletedByIDParams{
		ID:        waterID,
		DeletedBy: nullUserID,
	}
	err = cfg.db.MarkWaterNeedAsDeletedByID(r.Context(), deleteParams)
	if err != nil {
		cfg.sl.Debug("Could not mark water as deleted", "error", err, "water id", waterID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// respond with 204
	cfg.sl.Debug("Admin successfully marked plant water as deleted", "admin id", requestUserID, "water id", waterID)
	w.WriteHeader(http.StatusNoContent)
}

// TODO: ensure resource is sent back
// POST /admin/water/{water id} ? plant species id = uuid
func (cfg *apiConfig) adminSetPlantAsWaterNeedHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	waterIDStr := r.PathValue("waterID")
	waterID, err := uuid.Parse(waterIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse water id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		cfg.sl.Debug("No plant species id was provided in url query")
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
	nullWaterID := uuid.NullUUID{UUID: waterID, Valid: true}
	setParams := database.SetPlantSpeciesAsWaterNeedParams{
		ID:           plantSpeciesID,
		WaterNeedsID: nullWaterID,
		UpdatedBy:    requestUserID,
	}
	waterRecord, err := cfg.db.SetPlantSpeciesAsWaterNeed(r.Context(), setParams)
	if err != nil {
		cfg.sl.Debug("Could not set water for plant species", "error", err, "water id", waterID, "pant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	waterResponse := AdminSetWaterResponse{
		WaterNeedID:      waterRecord.ID,
		PlantSpeciesID:   plantSpeciesID,
		PlantSpeciesName: waterRecord.SpeciesName,
	}
	waterData, err := json.Marshal(&waterResponse)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully set plant species to water need", "admin id", requestUserID, "plant species id", plantSpeciesID, "water need", waterID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(waterData)
}

// DELETE /admin/water/{water id} ? plant species id = uuid
func (cfg *apiConfig) adminUnsetPlantAsWaterNeedHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	waterIDStr := r.PathValue("waterID")
	waterID, err := uuid.Parse(waterIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse water id from url path", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		cfg.sl.Debug("No plant species id was provided in url query")
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
	unsetParams := database.UnsetPlantSpeciesAsWaterNeedParams{
		ID:        plantSpeciesID,
		UpdatedBy: requestUserID,
	}
	waterRecord, err := cfg.db.UnsetPlantSpeciesAsWaterNeed(r.Context(), unsetParams)
	if err != nil {
		cfg.sl.Debug("Could not unset water need for plant species", "error", err, "water id", waterID, "plant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	waterResponse := AdminUnsetWaterResponse{
		PlantSpeciesID:   plantSpeciesID,
		PlantSpeciesName: waterRecord.SpeciesName,
	}
	waterData, err := json.Marshal(&waterResponse)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully unset plant species to water need", "admin id", requestUserID, "plant species id", plantSpeciesID, "water need", waterID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(waterData)
}
