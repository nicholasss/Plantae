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
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
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

func (cfg *apiConfig) adminPlantNamesViewHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		log.Printf("Could not get User ID from token due to: %q", err)
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	requestedLangCode := r.URL.Query().Get("lang")
	if requestedLangCode == "" {
		log.Print("Language filter not requested in URL query path.")
	} else {
		log.Printf("Language filter of %q requested in URL query path.", requestedLangCode)
	}

	// perform query without language filter
	if requestedLangCode == "" {
		plantNameRecords, err := cfg.db.GetAllPlantNamesOrderedByCreated(r.Context())
		if err != nil {
			log.Printf("Could not get plant name records without language code due to %q.", err)
			respondWithError(err, http.StatusInternalServerError, w)
			return
		}

		plantNameData, err := json.Marshal(plantNameRecords)
		if err != nil {
			log.Printf("Could not marshall plant name records due to %q.", err)
			respondWithError(err, http.StatusInternalServerError, w)
			return
		}

		log.Printf("Admin %q successfully queried plant names without language code.", requestUserID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(plantNameData)
		return
	}

	// perform query with language filter, checking code first
	requestedLangName, ok := LangCodes[requestedLangCode]
	if !ok {
		log.Printf("Unable to find language code of %q from query.", requestedLangCode)
		respondWithError(errors.New("language code requested does not exist"), http.StatusBadRequest, w)
		return
	}

	log.Printf("Language %q requested via LangCode of %q.", requestedLangName, requestedLangCode)
	plantNameRecords, err := cfg.db.GetAllPlantNamesForLanguageOrderedByCreated(r.Context(), requestedLangCode)
	if err != nil {
		log.Printf("Could not get plant name records with language code %q due to %q.", requestedLangCode, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	plantNameData, err := json.Marshal(plantNameRecords)
	if err != nil {
		log.Printf("Could not marshall plant name records due to %q.", err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully queried plant names with language code %q.", requestUserID, requestedLangCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(plantNameData)
}

func (cfg *apiConfig) adminPlantNamesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	plantNameIDStr := r.PathValue("plantNameID")
	plantNameID, err := uuid.Parse(plantNameIDStr)
	if err != nil {
		log.Printf("Could not parse plant name id from url path due to: %q", err)
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
	requestUserNullUUID := uuid.NullUUID{
		UUID:  requestUserID,
		Valid: true,
	}

	requestParams := database.MarkPlantNameAsDeletedByIDParams{
		ID:        plantNameID,
		DeletedBy: requestUserNullUUID,
	}
	err = cfg.db.MarkPlantNameAsDeletedByID(r.Context(), requestParams)
	if err != nil {
		log.Printf("Could not mark plant name %q as deleted due to: %q", plantNameID, err)
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Admin %q successfully marked plant name %q as deleted.", requestUserID, plantNameID)
	w.WriteHeader(http.StatusNoContent)
}
