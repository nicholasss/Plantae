package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nicholasss/plantae/internal/auth"
)

// === request response types

// requires that a admin access token be in the auth header
func (cfg *apiConfig) adminPlantsViewHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestAccessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	requestUserID, err := auth.ValidateJWT(requestAccessToken, cfg.JWTSecret)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), requestUserID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	if !userRecord.IsAdmin {
		respondWithError(fmt.Errorf("unauthorized"), http.StatusUnauthorized, w)
		return
	}
	// user is now authenticated below here

	plantSpeciesRecords, err := cfg.db.GetAllPlantSpeciesOrderedByCreated(r.Context())
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	if len(plantSpeciesRecords) <= 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	plantSpeciesData, err := json.Marshal(plantSpeciesRecords)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin listed all plant species successfully.")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(plantSpeciesData)
}
