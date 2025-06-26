package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
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

// === handler functions ===

// POST /admin/water
func (cfg *apiConfig) adminWaterCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var createRequest AdminWaterCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// checking body
	if createRequest.PlantType == "" {
		log.Print("Request Body missing plant type property.")
		respondWithError(errors.New("no plant type provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createRequest.Description == "" {
		log.Print("Request Body missing description property.")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	mmRequest := strings.ToLower(createRequest.PlantType) == "tropical" || strings.ToLower(createRequest.PlantType) == "temperate"
	dayRequest := strings.ToLower(createRequest.PlantType) == "semi-arid" || strings.ToLower(createRequest.PlantType) == "arid"

	var waterResponse AdminWaterResponse

	if dayRequest {
		if createRequest.DrySoilDays == nil {
			log.Print("Request Body missing dry soil days property.")
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
			log.Printf("Could not create water need record due to: %q", err)
			respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
			return
		}

		waterResponse.ID = waterRecord.ID
		waterResponse.PlantType = waterRecord.PlantType
		waterResponse.Description = waterRecord.Description
		waterResponse.DrySoilDays = &waterRecord.DrySoilDays.Int32

	} else if mmRequest {
		if createRequest.DrySoilMM == nil {
			log.Print("Request Body missing dry soil mm property.")
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
			log.Printf("Could not create water need record due to: %q", err)
			respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
			return
		}

		waterResponse.ID = waterRecord.ID
		waterResponse.PlantType = waterRecord.PlantType
		waterResponse.Description = waterRecord.Description
		waterResponse.DrySoilDays = &waterRecord.DrySoilMm.Int32

	} else {
		log.Printf("Invalid plant type of %q", createRequest.PlantType)
		respondWithError(errors.New("invalid plant type"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	waterData, err := json.Marshal(&waterResponse)
	if err != nil {
		log.Printf("Could not marshal records to json due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	log.Printf("Admin %q created water need successfully", requestUserID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(waterData)
}

func (cfg *apiConfig) adminWaterViewHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	waterRecords, err := cfg.db.GetAllWaterNeedsOrderedByCreated(r.Context())
	if err != nil {
		log.Printf("Could not get water need records due to: %q", err)
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
			log.Printf("Unknown plant type of %q in returned from query", record.PlantType)
		}
	}

	waterData, err := json.Marshal(waterResponses)
	if err != nil {
		log.Printf("Could not marshal records to json due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	log.Printf("Admin %q listed all water needs successfully", requestUserID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(waterData)
}

// DELETE /admin/water/{waterID}
func (cfg *apiConfig) adminWaterDeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// water need record to mark as deleted
	waterIDStr := r.PathValue("waterID")
	waterID, err := uuid.Parse(waterIDStr)
	if err != nil {
		log.Printf("Could not parse water id from url path due to: %q", err)
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
		log.Printf("Could not mark record as deleted due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// respond with 204
	log.Printf("Admin %q successfully deleted water need record %q", requestUserID, waterID)
	w.WriteHeader(http.StatusNoContent)
}

// POST /admin/water/{water id} ? plant species id = uuid
func (cfg *apiConfig) adminSetPlantAsWaterNeedHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	waterIDStr := r.PathValue("waterID")
	waterID, err := uuid.Parse(waterIDStr)
	if err != nil {
		log.Printf("Could not parse water need id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		log.Print("No plant species id was specified in url query")
		respondWithError(errors.New("no plant species id was provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		log.Printf("Could not parse plant species id from url query due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// user id
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
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
	err = cfg.db.SetPlantSpeciesAsWaterNeed(r.Context(), setParams)
	if err != nil {
		log.Printf("Could not set water needs %q for plant species %q due to: %q", waterID, plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	log.Printf("Admin %q successfully set water need %q for species %q ", requestUserID, waterID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /admin/water/{water id} ? plant species id = uuid
func (cfg *apiConfig) adminUnsetPlantAsWaterNeedHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	waterIDStr := r.PathValue("waterID")
	waterID, err := uuid.Parse(waterIDStr)
	if err != nil {
		log.Printf("Could not parse water need id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plant-species-id")
	if plantSpeciesIDStr == "" {
		log.Print("No plant species id was specified in url query")
		respondWithError(errors.New("no plant species id was provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	plantSpeciesID, err := uuid.Parse(plantSpeciesIDStr)
	if err != nil {
		log.Printf("Could not parse plant species id from url query due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// user id
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// perform unset
	unsetParams := database.UnsetPlantSpeciesAsWaterNeedParams{
		ID:        plantSpeciesID,
		UpdatedBy: requestUserID,
	}
	err = cfg.db.UnsetPlantSpeciesAsWaterNeed(r.Context(), unsetParams)
	if err != nil {
		log.Printf("Could not unset water need %q for plant species %q due to: %q", waterID, plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	log.Printf("Admin %q successfully unset water need %q for species %q ", requestUserID, waterID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}
