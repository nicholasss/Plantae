package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/database"
)

// There should only be four types of plants:
// - Tropical
// - Temperate
// - Semi-Arid
// - Arid

// === request response types ===

type AdminPlantTypeCreateRequest struct {
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	MaxTemperatureCelsius *int32  `json:"maxTemperatureCelsius"`
	MinTemperatureCelsius *int32  `json:"minTemperatureCelsius"`
	MaxHumidityPercent    *int32  `json:"maxHumidityPercent"`
	MinHumidityPercent    *int32  `json:"minHumidityPercent"`
	SoilOrganicMix        *string `json:"soilOrganicMix"`
	SoilGritMix           *string `json:"soilGritMix"`
	SoilDrainageMix       *string `json:"soilDrainageMix"`
}

type AdminPlantTypeCreateResponse struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	MaxTemperatureCelsius *int32    `json:"maxTemperatureCelsius,omitempty"`
	MinTemperatureCelsius *int32    `json:"minTemperatureCelsius,omitempty"`
	MaxHumidityPercent    *int32    `json:"maxHumidityPercent,omitempty"`
	MinHumidityPercent    *int32    `json:"minHumidityPercent,omitempty"`
	SoilOrganicMix        *string   `json:"soilOrganicMix,omitempty"`
	SoilGritMix           *string   `json:"soilGritMix,omitempty"`
	SoilDrainageMix       *string   `json:"soilDrainageMix,omitempty"`
}

type AdminPlantTypeViewResponse struct {
	ID                    uuid.UUID `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	MaxTemperatureCelsius *int32    `json:"maxTemperatureCelsius,omitempty"`
	MinTemperatureCelsius *int32    `json:"minTemperatureCelsius,omitempty"`
	MaxHumidityPercent    *int32    `json:"maxHumidityPercent,omitempty"`
	MinHumidityPercent    *int32    `json:"minHumidityPercent,omitempty"`
	SoilOrganicMix        *string   `json:"soilOrganicMix,omitempty"`
	SoilGritMix           *string   `json:"soilGritMix,omitempty"`
	SoilDrainageMix       *string   `json:"soilDrainageMix,omitempty"`
}

type AdminPlantTypeUpdateRequest struct {
	MaxTemperatureCelsius *int32  `json:"maxTemperatureCelsius"`
	MinTemperatureCelsius *int32  `json:"minTemperatureCelsius"`
	MaxHumidityPercent    *int32  `json:"maxHumidityPercent"`
	MinHumidityPercent    *int32  `json:"minHumidityPercent"`
	SoilOrganicMix        *string `json:"soilOrganicMix"`
	SoilGritMix           *string `json:"soilGritMix"`
	SoilDrainageMix       *string `json:"soilDrainageMix"`
}

type AdminSetPlantTypeResponse struct {
	PlantTypeID      uuid.UUID `json:"plantTypeID"`
	PlantSpeciesID   uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName string    `json:"plantSpeciesName"`
}

type AdminUnsetPlantTypeResponse struct {
	PlantSpeciesID   uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName string    `json:"plantSpeciesName"`
}

// === handler functions ===

// POST /api/v1/admin/plant-types
// Create plant type request
func (cfg *apiConfig) adminPlantTypesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var createRequest AdminPlantTypeCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// checking body properties
	if createRequest.Name == "" {
		cfg.sl.Debug("Request body missing name")
		respondWithError(errors.New("no name provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createRequest.Description == "" {
		cfg.sl.Debug("Request body missing description")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	// converting properties
	maxTC := sql.NullInt32{}
	minTC := sql.NullInt32{}
	maxHP := sql.NullInt32{}
	minHP := sql.NullInt32{}
	soilOM := sql.NullString{}
	soilGM := sql.NullString{}
	soilDM := sql.NullString{}

	if createRequest.MaxTemperatureCelsius == nil {
		maxTC.Valid = false
	} else {
		maxTC.Valid = true
		maxTC.Int32 = *createRequest.MaxTemperatureCelsius
	}
	if createRequest.MinTemperatureCelsius == nil {
		minTC.Valid = false
	} else {
		minTC.Valid = true
		minTC.Int32 = *createRequest.MinTemperatureCelsius
	}
	if createRequest.MaxHumidityPercent == nil {
		maxHP.Valid = false
	} else {
		maxHP.Valid = true
		maxHP.Int32 = *createRequest.MaxHumidityPercent
	}
	if createRequest.MinHumidityPercent == nil {
		minHP.Valid = false
	} else {
		minHP.Valid = true
		minHP.Int32 = *createRequest.MinHumidityPercent
	}
	if createRequest.SoilOrganicMix == nil {
		soilOM.Valid = false
	} else {
		soilOM.Valid = true
		soilOM.String = *createRequest.SoilOrganicMix
	}
	if createRequest.SoilGritMix == nil {
		soilGM.Valid = false
	} else {
		soilGM.Valid = true
		soilGM.String = *createRequest.SoilGritMix
	}
	if createRequest.SoilDrainageMix == nil {
		soilDM.Valid = false
	} else {
		soilDM.Valid = true
		soilDM.String = *createRequest.SoilDrainageMix
	}

	createParams := database.CreatePlantTypeParams{
		CreatedBy:             requestUserID,
		Name:                  createRequest.Name,
		Description:           createRequest.Description,
		MaxTemperatureCelsius: maxTC,
		MinTemperatureCelsius: minTC,
		MaxHumidityPercent:    maxHP,
		MinHumidityPercent:    minHP,
		SoilOrganicMix:        soilOM,
		SoilGritMix:           soilGM,
		SoilDrainageMix:       soilDM,
	}
	typeRecord, err := cfg.db.CreatePlantType(r.Context(), createParams)
	if err != nil {
		cfg.sl.Debug("Could not create plant type in database", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	var MaxTemperatureCelsius *int32
	var MinTemperatureCelsius *int32
	var MaxHumidityPercent *int32
	var MinHumidityPercent *int32
	var SoilOrganicMix *string
	var SoilGritMix *string
	var SoilDrainageMix *string

	if typeRecord.MaxTemperatureCelsius.Valid {
		MaxTemperatureCelsius = &typeRecord.MaxTemperatureCelsius.Int32
	}
	if typeRecord.MinTemperatureCelsius.Valid {
		MinTemperatureCelsius = &typeRecord.MinTemperatureCelsius.Int32
	}
	if typeRecord.MaxHumidityPercent.Valid {
		MaxHumidityPercent = &typeRecord.MaxHumidityPercent.Int32
	}
	if typeRecord.MinHumidityPercent.Valid {
		MinHumidityPercent = &typeRecord.MinHumidityPercent.Int32
	}
	if typeRecord.SoilOrganicMix.Valid {
		SoilOrganicMix = &typeRecord.SoilOrganicMix.String
	}
	if typeRecord.SoilGritMix.Valid {
		SoilGritMix = &typeRecord.SoilGritMix.String
	}
	if typeRecord.SoilDrainageMix.Valid {
		SoilDrainageMix = &typeRecord.SoilDrainageMix.String
	}

	plantTypeResponse := AdminPlantTypeViewResponse{
		ID:                    typeRecord.ID,
		Name:                  typeRecord.Name,
		Description:           typeRecord.Description,
		MaxTemperatureCelsius: MaxTemperatureCelsius,
		MinTemperatureCelsius: MinTemperatureCelsius,
		MaxHumidityPercent:    MaxHumidityPercent,
		MinHumidityPercent:    MinHumidityPercent,
		SoilOrganicMix:        SoilOrganicMix,
		SoilGritMix:           SoilGritMix,
		SoilDrainageMix:       SoilDrainageMix,
	}

	cfg.sl.Debug("Admin successfully create plant type", "admin id", requestUserID, "plant type id", typeRecord.ID)
	respondWithJSON(http.StatusCreated, plantTypeResponse, w, cfg.sl)
}

// GET /admin/plant-type
// view list of plant types
func (cfg *apiConfig) adminPlantTypesViewHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	plantTypeRecords, err := cfg.db.GetAllPlantTypesOrderedByCreated(r.Context())
	if err != nil {
		cfg.sl.Debug("Could not get plant type records", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	if len(plantTypeRecords) <= 0 {
		cfg.sl.Debug("Admin successfully listed empty plant species list", "admin id", requestUserID)
		respondWithJSON(http.StatusOK, plantTypeRecords, w, cfg.sl)
		return
	}

	// TODO: not the most efficient way to convert, is there another way?
	plantTypeResponse := make([]AdminPlantTypeViewResponse, 0)
	for _, oldRecord := range plantTypeRecords {
		var MaxTemperatureCelsius *int32
		var MinTemperatureCelsius *int32
		var MaxHumidityPercent *int32
		var MinHumidityPercent *int32
		var SoilOrganicMix *string
		var SoilGritMix *string
		var SoilDrainageMix *string

		if oldRecord.MaxTemperatureCelsius.Valid {
			MaxTemperatureCelsius = &oldRecord.MaxTemperatureCelsius.Int32
		}
		if oldRecord.MinTemperatureCelsius.Valid {
			MinTemperatureCelsius = &oldRecord.MinTemperatureCelsius.Int32
		}
		if oldRecord.MaxHumidityPercent.Valid {
			MaxHumidityPercent = &oldRecord.MaxHumidityPercent.Int32
		}
		if oldRecord.MinHumidityPercent.Valid {
			MinHumidityPercent = &oldRecord.MinHumidityPercent.Int32
		}
		if oldRecord.SoilOrganicMix.Valid {
			SoilOrganicMix = &oldRecord.SoilOrganicMix.String
		}
		if oldRecord.SoilGritMix.Valid {
			SoilGritMix = &oldRecord.SoilGritMix.String
		}
		if oldRecord.SoilDrainageMix.Valid {
			SoilDrainageMix = &oldRecord.SoilDrainageMix.String
		}

		newRecord := AdminPlantTypeViewResponse{
			ID:                    oldRecord.ID,
			Name:                  oldRecord.Name,
			Description:           oldRecord.Description,
			MaxTemperatureCelsius: MaxTemperatureCelsius,
			MinTemperatureCelsius: MinTemperatureCelsius,
			MaxHumidityPercent:    MaxHumidityPercent,
			MinHumidityPercent:    MinHumidityPercent,
			SoilOrganicMix:        SoilOrganicMix,
			SoilGritMix:           SoilGritMix,
			SoilDrainageMix:       SoilDrainageMix,
		}

		plantTypeResponse = append(plantTypeResponse, newRecord)
	}

	cfg.sl.Debug("Admin successfully listed plant type list", "admin id", requestUserID)
	respondWithJSON(http.StatusOK, plantTypeResponse, w, cfg.sl)
}

// plant type info update
func (cfg *apiConfig) adminPlantTypesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant type id from url path", "error", err)
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

	var updateRequest AdminPlantTypeUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// converting properties
	maxTC := sql.NullInt32{}
	minTC := sql.NullInt32{}
	maxHP := sql.NullInt32{}
	minHP := sql.NullInt32{}
	soilOM := sql.NullString{}
	soilGM := sql.NullString{}
	soilDM := sql.NullString{}

	if updateRequest.MaxTemperatureCelsius == nil {
		maxTC.Valid = false
	} else {
		maxTC.Valid = true
		maxTC.Int32 = *updateRequest.MaxTemperatureCelsius
	}
	if updateRequest.MinTemperatureCelsius == nil {
		minTC.Valid = false
	} else {
		minTC.Valid = true
		minTC.Int32 = *updateRequest.MinTemperatureCelsius
	}
	if updateRequest.MaxHumidityPercent == nil {
		maxHP.Valid = false
	} else {
		maxHP.Valid = true
		maxHP.Int32 = *updateRequest.MaxHumidityPercent
	}
	if updateRequest.MinHumidityPercent == nil {
		minHP.Valid = false
	} else {
		minHP.Valid = true
		minHP.Int32 = *updateRequest.MinHumidityPercent
	}
	if updateRequest.SoilOrganicMix == nil {
		soilOM.Valid = false
	} else {
		soilOM.Valid = true
		soilOM.String = *updateRequest.SoilOrganicMix
	}
	if updateRequest.SoilGritMix == nil {
		soilGM.Valid = false
	} else {
		soilGM.Valid = true
		soilGM.String = *updateRequest.SoilGritMix
	}
	if updateRequest.SoilDrainageMix == nil {
		soilDM.Valid = false
	} else {
		soilDM.Valid = true
		soilDM.String = *updateRequest.SoilDrainageMix
	}

	updateParams := database.UpdatePlantTypesPropertiesByIDParams{
		ID:                    plantTypeID,
		MaxTemperatureCelsius: maxTC,
		MinTemperatureCelsius: minTC,
		MaxHumidityPercent:    maxHP,
		MinHumidityPercent:    minHP,
		SoilOrganicMix:        soilOM,
		SoilGritMix:           soilGM,
		SoilDrainageMix:       soilDM,
	}

	err = cfg.db.UpdatePlantTypesPropertiesByID(r.Context(), updateParams)
	if err != nil {
		cfg.sl.Debug("Could not update plant type record", "error", err, "plant type id", plantTypeID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully updated plant type", "admin id", requestUserID, "plant type id", plantTypeID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/admin/plant-type/{plantTypeID}
func (cfg *apiConfig) adminPlantTypeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant type id from url path", "error", err)
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

	// perform delete
	nullAdminID := uuid.NullUUID{Valid: true, UUID: requestUserID}
	deleteParams := database.MarkPlantTypeAsDeletedByIDParams{
		ID:        plantTypeID,
		DeletedBy: nullAdminID,
	}
	err = cfg.db.MarkPlantTypeAsDeletedByID(r.Context(), deleteParams)
	if err != nil {
		cfg.sl.Debug("Could not mark plant type as deleted", "error", err, "plant type id", plantTypeID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully marked plant type as deleted", "admin id", requestUserID, "plant type id", plantTypeID)
	w.WriteHeader(http.StatusNoContent)
}

// POST /admin/plant-type/{plant type id} ? plant species id = uuid
func (cfg *apiConfig) adminSetPlantAsTypeHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant type id from url path", "error", err)
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
	nullPlantTypeID := uuid.NullUUID{Valid: true, UUID: plantTypeID}
	setPlantTypeParams := database.SetPlantSpeciesAsTypeParams{
		ID:          plantSpeciesID,
		PlantTypeID: nullPlantTypeID,
		UpdatedBy:   requestUserID,
	}
	speciesRecord, err := cfg.db.SetPlantSpeciesAsType(r.Context(), setPlantTypeParams)
	if err != nil {
		cfg.sl.Debug("Could not set plant type for plant species", "error", err, "plant type id", plantTypeID, "plant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	setResponse := AdminSetPlantTypeResponse{
		PlantTypeID:      plantTypeID,
		PlantSpeciesID:   plantSpeciesID,
		PlantSpeciesName: speciesRecord.SpeciesName,
	}

	cfg.sl.Debug("Admin successfully set plant species to plant type", "admin id", requestUserID, "plant species id", plantSpeciesID, "plant type id", plantTypeID)
	respondWithJSON(http.StatusOK, setResponse, w, cfg.sl)
}

// DELETE /admin/plant-type/{plant type id} ? plant species id = uuid
func (cfg *apiConfig) adminUnsetPlantAsTypeHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant type id from url path", "error", err)
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

	// perform operation
	unsetPlantTypeParams := database.UnsetPlantSpeciesAsTypeParams{
		ID:        plantSpeciesID,
		UpdatedBy: requestUserID,
	}
	speciesRecord, err := cfg.db.UnsetPlantSpeciesAsType(r.Context(), unsetPlantTypeParams)
	if err != nil {
		cfg.sl.Debug("Could not unset type for plant species", "error", err, "plant type id", plantTypeID, "plant species id", plantSpeciesID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	unsetResponse := AdminUnsetPlantTypeResponse{
		PlantSpeciesID:   plantSpeciesID,
		PlantSpeciesName: speciesRecord.SpeciesName,
	}

	cfg.sl.Debug("Admin successfully unset plant species from plant type", "admin id", requestUserID, "plant species id", plantSpeciesID, "plant type id", plantTypeID)
	respondWithJSON(http.StatusOK, unsetResponse, w, cfg.sl)
}
