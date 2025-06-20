package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

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
	PlantType   string `json:"plantType"`
	Description string `json:"description"`
	DrySoilMM   *int32 `json:"drySoilMM,omitempty"`
	DrySoilDays *int32 `json:"drySoilDays,omitempty"`
}

// === handler functions ===

// POST /admin/water
func (cfg *apiConfig) adminWaterCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var createRequest AdminWaterCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// checking body
	if createRequest.PlantType == "" {
		log.Print("Request Body missing plant type property.")
		respondWithError(errors.New("no plant type provided"), http.StatusBadRequest, w)
		return
	}
	if createRequest.Description == "" {
		log.Print("Request Body missing description property.")
		respondWithError(errors.New("no description provided"), http.StatusBadRequest, w)
		return
	}

	mmRequest := strings.ToLower(createRequest.PlantType) == "tropical" || strings.ToLower(createRequest.PlantType) == "temperate"
	dayRequest := strings.ToLower(createRequest.PlantType) == "semi-arid" || strings.ToLower(createRequest.PlantType) == "arid"

	if dayRequest {
		if createRequest.DrySoilDays == nil {
			log.Print("Request Body missing dry soil days property.")
			respondWithError(errors.New("no dry soil days provided"), http.StatusBadRequest, w)
			return
		}

		nullDrySoilDays := sql.NullInt32{Int32: *createRequest.DrySoilDays, Valid: true}
		createParams := database.CreateWaterDryDaysParams{
			CreatedBy:   requestUserID,
			PlantType:   createRequest.PlantType,
			Description: createRequest.Description,
			DrySoilDays: nullDrySoilDays,
		}
		_, err = cfg.db.CreateWaterDryDays(r.Context(), createParams)
		if err != nil {
			log.Printf("Could not create water need record due to: %q", err)
			respondWithError(err, http.StatusInternalServerError, w)
			return
		}

	} else if mmRequest {
		if createRequest.DrySoilMM == nil {
			log.Print("Request Body missing dry soil mm property.")
			respondWithError(errors.New("no dry soil mm provided"), http.StatusBadRequest, w)
			return
		}

		nullDrySoilMM := sql.NullInt32{Int32: *createRequest.DrySoilMM, Valid: true}
		createParams := database.CreateWaterDryMMParams{
			CreatedBy:   requestUserID,
			PlantType:   createRequest.PlantType,
			Description: createRequest.Description,
			DrySoilMm:   nullDrySoilMM,
		}
		_, err = cfg.db.CreateWaterDryMM(r.Context(), createParams)
		if err != nil {
			log.Printf("Could not create water need record due to: %q", err)
			respondWithError(err, http.StatusInternalServerError, w)
			return
		}

	} else {
		log.Printf("Invalid plant type of %q", createRequest.PlantType)
		respondWithError(errors.New("invalid plant type"), http.StatusBadRequest, w)
		return
	}

	log.Printf("Admin %q created water need successfully", requestUserID)
	w.WriteHeader(http.StatusNoContent)
}
