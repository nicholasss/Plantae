package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
)

type UserViewAllPlantInfoResponse struct {
	PlantSpeciesID       uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName     string    `json:"plantSpeciesName"`
	HumanPoisonToxic     *bool     `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic       *bool     `json:"petPoisonToxic,omitempty"`
	HumanEdible          *bool     `json:"humanEdible,omitempty"`
	PetEdible            *bool     `json:"petEdible,omitempty"`
	PlantTypeName        *string   `json:"plantTypeName,omitempty"`
	PlantTypeDescription *string   `json:"plantTypeDescription,omitempty"`
	LightNeedName        *string   `json:"lightNeedName,omitempty"`
	LightNeedDescription *string   `json:"lightNeedDescription,omitempty"`
	WaterNeedName        *string   `json:"waterNeedName,omitempty"`
	WaterNeedDescription *string   `json:"waterNeedDescription,omitmepty"`
	WaterNeedDrySoilMM   *int32    `json:"waterNeedDrySoilMM,omitempty"`
	WaterNeedDrySoilDays *int32    `json:"waterNeedDrySoilDays,omitempty"`
}

func (cfg *apiConfig) usersViewPlantsListHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenProvided, err := auth.GetBearerToken(r.Header, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get token from headers", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	requestUserID, err := auth.ValidateJWT(accessTokenProvided, cfg.JWTSecret, cfg.sl)
	if err != nil {
		cfg.sl.Debug("Could not get user id from token", "error", err)
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	requestedLangCode := r.URL.Query().Get("lang-code")
	if requestedLangCode == "" {
		cfg.sl.Debug("No lang code was provided in query params")
		respondWithError(errors.New("no lang code was requested in query params"), http.StatusBadRequest, w, cfg.sl)
		return
	}
	langName, ok := LangCodes[requestedLangCode]
	if !ok {
		cfg.sl.Debug("Unknown language is being requested", "lang code", requestedLangCode)
		respondWithError(errors.New("invalid lang code was requested"), http.StatusBadRequest, w, cfg.sl)
	}
	cfg.sl.Debug("Searching for all plants with common names in language", "lang code", requestedLangCode, "lang name", langName)

	plantRecords, err := cfg.db.GetAllViewPlantsOrderedByUpdated(r.Context(), requestedLangCode)
	if err != nil {
		cfg.sl.Debug("Could not view all plants in database with lang code", "error", err, "lang code", requestedLangCode)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	if len(plantRecords) <= 0 {
		cfg.sl.Debug("User successfully viewed empty plants view list", "user id", requestUserID)
		respondWithJSON(http.StatusOK, plantRecords, w, cfg.sl)
		return
	}

	// conversion to response records
	var plantResponses []UserViewAllPlantInfoResponse
	for _, record := range plantRecords {
		mmRecord := strings.ToLower(record.WaterNeedType) == "tropical" || strings.ToLower(record.WaterNeedType) == "temperate"
		dayRecord := strings.ToLower(record.WaterNeedType) == "semi-arid" || strings.ToLower(record.WaterNeedType) == "arid"
	}

	cfg.sl.Debug("User successfully listed all available plants", "user id", requestUserID)
}
