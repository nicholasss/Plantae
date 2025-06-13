package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types ===

// create request for plant name
type AdminPlantNamesCreateRequest struct {
	PlantID    uuid.UUID `json:"plantID"`
	LangCode   string    `json:"langCode"`
	CommonName string    `json:"commonName"`
}

func (cfg *apiConfig) adminPlantNamesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.authorizeNormalAdmin(r)
	if err != nil {
		log.Printf("Could not authorize normal (non-superadmin) due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	var createRequest AdminPlantNamesCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		log.Printf("Could not decode body of request due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// check all request properties
	if createRequest.PlantID == uuid.Nil {
		log.Print("Request body missing plant id.")
		respondWithError(errors.New("no plant id provided"), http.StatusBadRequest, w)
		return
	}
	if createRequest.LangCode == "" {
		log.Print("Request body missing lang code.")
		respondWithError(errors.New("no lang code provided"), http.StatusBadRequest, w)
		return
	}
	if createRequest.CommonName == "" {
		log.Print("Request body missing common name.")
		respondWithError(errors.New("no common name provided"), http.StatusBadRequest, w)
		return
	}

	createRequestParams := database.CreatePlantNameParams{
		CreatedBy:  requestUserID,
		PlantID:    createRequest.PlantID,
		LangCode:   createRequest.LangCode,
		CommonName: createRequest.CommonName,
	}

	_, err = cfg.db.CreatePlantName(r.Context(), createRequestParams)
	if err != nil {
		log.Printf("Could not create plant name record for plant id %q due to: %q", createRequest.PlantID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q created plant name record %q for plant id %q", requestUserID, createRequest.CommonName, createRequest.PlantID)
	w.WriteHeader(http.StatusNoContent)
}
