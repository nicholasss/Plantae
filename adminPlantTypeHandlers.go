package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

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

type AdminPlantTypeViewResponse struct {
	ID                    uuid.UUID `json:"id"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	CreatedBy             uuid.UUID `json:"createdBy"`
	UpdatedBy             uuid.UUID `json:"updatedBy"`
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

// === handler functions ===

// POST /admin/plant-type
// Create plant type request
func (cfg *apiConfig) adminPlantTypesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var createRequest AdminPlantTypeCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// checking body properties
	if createRequest.Name == "" {
		log.Print("Request Body missing name proptery.")
		respondWithError(errors.New("no name provided"), http.StatusBadRequest, w)
		return
	}
	if createRequest.Description == "" {
		log.Print("Request body missing description property.")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w)
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

	_, err = cfg.db.CreatePlantType(r.Context(), createParams)
	if err != nil {
		log.Printf("Could not create plant type record due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// created successfully
	log.Printf("Admin %q created plant type record successfully.", requestUserID)
	w.WriteHeader(http.StatusNoContent)
}

// GET /admin/plant-type
// view list of plant types
func (cfg *apiConfig) adminPlantTypesViewHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	plantTypeRecords, err := cfg.db.GetAllPlantTypesOrderedByCreated(r.Context())
	if err != nil {
		log.Printf("Could not get plant type records due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	if len(plantTypeRecords) <= 0 {
		log.Print("Admin listed empty plant types list successfully.")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	plantTypesResponse := make([]AdminPlantTypeViewResponse, 0)
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
			CreatedAt:             oldRecord.CreatedAt,
			UpdatedAt:             oldRecord.UpdatedAt,
			CreatedBy:             oldRecord.CreatedBy,
			UpdatedBy:             oldRecord.UpdatedBy,
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

		plantTypesResponse = append(plantTypesResponse, newRecord)
	}

	plantTypesData, err := json.Marshal(plantTypesResponse)
	if err != nil {
		log.Printf("Could not marshal records to json due to: %q", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q isted plant types list successfully.", requestUserID)
	// log.Printf("DEBUG: list of plants types: %s", string(plantTypesData))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(plantTypesData)
}

// plant type info update
func (cfg *apiConfig) adminPlantTypesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		log.Printf("Could not parse plant type id from url path due to: %q", err)
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

	var updateRequest AdminPlantTypeUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
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
		log.Printf("Could not update plant type record %q due to: %q", plantTypeID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q updated plant type %q successfully.", requestUserID, plantTypeID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/v1/admin/plant-type/{plantTypeID}
func (cfg *apiConfig) adminPlantTypeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		log.Printf("Could not parse plant type id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
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
		log.Printf("Could not mark plant type %q as deleted due to: %q", plantTypeID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully marked plant type %q as deleted.", requestUserID, plantTypeID)
	w.WriteHeader(http.StatusNoContent)
}

// POST /admin/plant-type/{plant type id} ? plant species id = uuid
func (cfg *apiConfig) adminSetPlantAsTypeHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		log.Printf("Could not parse plant type id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plantSpeciesID")
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
	nullPlantTypeID := uuid.NullUUID{Valid: true, UUID: plantTypeID}
	setPlantTypeParams := database.SetPlantSpeciesAsTypeParams{
		ID:          plantSpeciesID,
		PlantTypeID: nullPlantTypeID,
		UpdatedBy:   requestUserID,
	}
	err = cfg.db.SetPlantSpeciesAsType(r.Context(), setPlantTypeParams)
	if err != nil {
		log.Printf("Could not set type %q for plant species %q due to: %q", plantTypeID, plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully set type %q for species %q ", requestUserID, plantTypeID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /admin/plant-type/{plant type id} ? plant species id = uuid
func (cfg *apiConfig) adminUnsetPlantAsTypeHandler(w http.ResponseWriter, r *http.Request) {
	// plant type
	plantTypeIDStr := r.PathValue("plantTypeID")
	plantTypeID, err := uuid.Parse(plantTypeIDStr)
	if err != nil {
		log.Printf("Could not parse plant type id from url path due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// plant species
	plantSpeciesIDStr := r.URL.Query().Get("plantSpeciesID")
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

	// perform operation
	unsetPlantTypeParams := database.UnsetPlantSpeciesAsTypeParams{
		ID:        plantSpeciesID,
		UpdatedBy: requestUserID,
	}
	err = cfg.db.UnsetPlantSpeciesAsType(r.Context(), unsetPlantTypeParams)
	if err != nil {
		log.Printf("Could not unset type %q for plant species %q due to: %q", plantTypeID, plantSpeciesID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully unset type %q for species %q ", requestUserID, plantTypeID, plantSpeciesID)
	w.WriteHeader(http.StatusNoContent)
}
