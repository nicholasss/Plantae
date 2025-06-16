package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nicholasss/plantae/internal/database"
)

// There should only be four types of plants:
// - Tropical
// - Temperate
// - Semi-Arid
// - Arid

// === request response types ===

type AdminPlantTypeCreateRequest struct {
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	MaxTemperatureCelsius *float64 `json:"maxTemperatureCelsius"`
	MinTemperatureCelsius *float64 `json:"minTemperatureCelsius"`
	MaxHumidityPercent    *float64 `json:"maxHumidityPercent"`
	MinHumidityPercent    *float64 `json:"minHumidityPercent"`
	SoilOrganicMix        *string  `json:"soilOrganicMix"`
	SoilGritMix           *string  `json:"soilGritMix"`
	SoilDrainageMix       *string  `json:"soilDrainageMix"`
}

// === handler functions ===

// POST /admin/plant-type
// Create plant type request
func (cfg *apiConfig) adminPlantTypesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.authorizeNormalAdmin(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
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
	maxTC := sql.NullFloat64{}
	minTC := sql.NullFloat64{}
	maxHP := sql.NullFloat64{}
	minHP := sql.NullFloat64{}
	soilOM := sql.NullString{}
	soilGM := sql.NullString{}
	soilDM := sql.NullString{}

	if createRequest.MaxTemperatureCelsius == nil {
		maxTC.Valid = false
	} else {
		maxTC.Valid = true
		maxTC.Float64 = *createRequest.MaxTemperatureCelsius
	}
	if createRequest.MinTemperatureCelsius == nil {
		minTC.Valid = false
	} else {
		minTC.Valid = true
		minTC.Float64 = *createRequest.MinTemperatureCelsius
	}
	if createRequest.MaxHumidityPercent == nil {
		maxHP.Valid = false
	} else {
		maxHP.Valid = true
		maxHP.Float64 = *createRequest.MaxHumidityPercent
	}
	if createRequest.MinHumidityPercent == nil {
		minHP.Valid = false
	} else {
		minHP.Valid = true
		minHP.Float64 = *createRequest.MinHumidityPercent
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
