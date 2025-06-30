package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/database"
)

// === request response types ===

// AdminPlantNamesCreateRequest is for decoding plant name requests.
type AdminPlantNamesCreateRequest struct {
	PlantID    uuid.UUID `json:"plantID"`
	LangCode   string    `json:"langCode"`
	CommonName string    `json:"commonName"`
}

// AdminPlantNamesCreateResponse is for encoding plant name responses.
type AdminPlantNamesCreateResponse struct {
	ID         uuid.UUID `json:"ID"`
	PlantID    uuid.UUID `json:"plantID"`
	LangCode   string    `json:"langCode"`
	CommonName string    `json:"commonName"`
}

// TODO: ensure resource is sent back
// POST /api/v1/admin/plant-names
func (cfg *apiConfig) adminPlantNamesCreateHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	var createRequest AdminPlantNamesCreateRequest
	err = json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		cfg.sl.Debug("Could not decode body of request", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// check all request properties
	if createRequest.PlantID == uuid.Nil {
		cfg.sl.Debug("Request body missing plant id")
		respondWithError(errors.New("no plant id provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createRequest.LangCode == "" {
		cfg.sl.Debug("Request body missing lang code")
		respondWithError(errors.New("no lang code provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	if createRequest.CommonName == "" {
		cfg.sl.Debug("Request body missing common name")
		respondWithError(errors.New("no common name provided"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	createRequestParams := database.CreatePlantNameParams{
		CreatedBy:  requestUserID,
		PlantID:    createRequest.PlantID,
		LangCode:   createRequest.LangCode,
		CommonName: createRequest.CommonName,
	}

	plantNameRecord, err := cfg.db.CreatePlantName(r.Context(), createRequestParams)
	if err != nil {
		cfg.sl.Debug("Could not create plant name record for plant id", "error", err, "plant id", createRequest.PlantID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	nameData, err := json.Marshal(&plantNameRecord)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin created plant name record", "admin id", requestUserID, "common name", createRequest.CommonName, "plant id", createRequest.PlantID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(nameData)
}

func (cfg *apiConfig) adminPlantNamesViewHandler(w http.ResponseWriter, r *http.Request) {
	// check header for admin access token
	requestUserID, err := cfg.getUserIDFromToken(r)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	requestedLangCode := r.URL.Query().Get("lang")
	if requestedLangCode == "" {
		cfg.sl.Debug("Language filter not requested in URL query path")
	} else {
		cfg.sl.Debug("Language filter requested in URL query path", "lang code", requestedLangCode)
	}

	// perform query without language filter
	if requestedLangCode == "" {
		plantNameRecords, err := cfg.db.GetAllPlantNamesOrderedByCreated(r.Context())
		if err != nil {
			cfg.sl.Debug("Could not get plant name records without lang code", "error", err)
			respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
			return
		}

		plantNameData, err := json.Marshal(plantNameRecords)
		if err != nil {
			cfg.sl.Debug("Could not marshal data", "error", err)
			respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
			return
		}

		cfg.sl.Debug("Admin successfully queried all common names", "admin id", requestUserID)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(plantNameData)
		return
	}

	// perform query with language filter, checking code first
	requestedLangName, ok := LangCodes[requestedLangCode]
	if !ok {
		cfg.sl.Debug("Requested language code was not found", "lang code", requestedLangCode)
		respondWithError(errors.New("language code requested does not exist"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Filtering common names to show requested lang code", "lang code", requestedLangCode, "lang name", requestedLangName)
	plantNameRecords, err := cfg.db.GetAllPlantNamesForLanguageOrderedByCreated(r.Context(), requestedLangCode)
	if err != nil {
		cfg.sl.Debug("Could not get common names for language code", "error", err, "lang code", requestedLangCode)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	plantNameData, err := json.Marshal(plantNameRecords)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin successfully queried common names for language", "admin id", requestUserID, "lang code", requestedLangCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(plantNameData)
}

func (cfg *apiConfig) adminPlantNamesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	plantNameIDStr := r.PathValue("plantNameID")
	plantNameID, err := uuid.Parse(plantNameIDStr)
	if err != nil {
		cfg.sl.Debug("Could not parse plant name id from url path", "error", err)
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
		cfg.sl.Debug("Could not mark plant name as deleted", "error", err, "plant name id", plantNameID)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Debug("Admin marked plant name record as deleted", "admin id", requestUserID, "plant name id", plantNameID)
	w.WriteHeader(http.StatusNoContent)
}
