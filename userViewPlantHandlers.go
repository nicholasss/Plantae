package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/nicholasss/plantae/internal/auth"
)

type UserViewAllPlantInfoResponse struct {
	PlantSpeciesID       uuid.UUID `json:"plantSpeciesID"`
	PlantSpeciesName     string    `json:"plantSpeciesName"`
	CommonNamesLangCode  *string   `json:"commonNameLangCode"`
	CommonNames          *string   `json:"commonNames"`
	HumanPoisonToxic     *bool     `json:"humanPoisonToxic,omitempty"`
	PetPoisonToxic       *bool     `json:"petPoisonToxic,omitempty"`
	HumanEdible          *bool     `json:"humanEdible,omitempty"`
	PetEdible            *bool     `json:"petEdible,omitempty"`
	PlantTypeName        *string   `json:"plantTypeName,omitempty"`
	PlantTypeDescription *string   `json:"plantTypeDescription,omitempty"`
	LightNeedName        *string   `json:"lightNeedName,omitempty"`
	LightNeedDescription *string   `json:"lightNeedDescription,omitempty"`
	WaterNeedName        *string   `json:"waterNeedName,omitempty"`
	WaterNeedDescription *string   `json:"waterNeedDescription,omitempty"`
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

	// NOTE: if this section is providing errors, set the internal/database package to use sql.NullString
	nullLangCode := sql.NullString{String: requestedLangCode, Valid: true}
	plantRecords, err := cfg.db.GetAllViewPlantsOrderedByUpdated(r.Context(), nullLangCode)
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

		var commonNamesLangCode *string
		var commonNames *string
		var humanPT *bool
		var petPT *bool
		var humanE *bool
		var petE *bool
		var plantTypeName *string
		var plantTypeDesc *string
		var lightNeedName *string
		var lightNeedDesc *string
		var waterNeedName *string
		var waterNeedDesc *string
		var waterNeedDryMM *int32
		var waterNeedDryDays *int32

		// NOTE: if this section is providing errors, set the internal/database package to use sql.NullString
		if record.LangCode.Valid {
			commonNamesLangCode = &record.LangCode.String
		}
		if record.CommonNames.Valid {
			commonNames = &record.CommonNames.String
		}

		if record.HumanPoisonToxic.Valid {
			humanPT = &record.HumanPoisonToxic.Bool
		}
		if record.PetPoisonToxic.Valid {
			petPT = &record.PetPoisonToxic.Bool
		}
		if record.HumanEdible.Valid {
			humanE = &record.HumanEdible.Bool
		}
		if record.PetEdible.Valid {
			petE = &record.PetEdible.Bool
		}
		if record.PlantTypeName.Valid {
			plantTypeName = &record.PlantTypeName.String
		}
		if record.PlantTypeDescription.Valid {
			plantTypeDesc = &record.PlantTypeDescription.String
		}
		if record.LightNeedName.Valid {
			lightNeedName = &record.LightNeedName.String
		}
		if record.LightNeedDescription.Valid {
			lightNeedDesc = &record.LightNeedDescription.String
		}
		if record.WaterNeedType.Valid {
			waterNeedName = &record.WaterNeedType.String
		}
		if record.WaterNeedDescription.Valid {
			waterNeedDesc = &record.WaterNeedDescription.String
		}
		if record.WaterNeedDrySoilMm.Valid {
			waterNeedDryMM = &record.WaterNeedDrySoilMm.Int32
		}
		if record.WaterNeedDrySoilDays.Valid {
			waterNeedDryDays = &record.WaterNeedDrySoilDays.Int32
		}

		response := UserViewAllPlantInfoResponse{
			PlantSpeciesID:       record.PlantSpeciesID,
			PlantSpeciesName:     record.PlantSpeciesName,
			CommonNamesLangCode:  commonNamesLangCode,
			CommonNames:          commonNames,
			HumanPoisonToxic:     humanPT,
			PetPoisonToxic:       petPT,
			HumanEdible:          humanE,
			PetEdible:            petE,
			PlantTypeName:        plantTypeName,
			PlantTypeDescription: plantTypeDesc,
			LightNeedName:        lightNeedName,
			LightNeedDescription: lightNeedDesc,
			WaterNeedName:        waterNeedName,
			WaterNeedDescription: waterNeedDesc,
			WaterNeedDrySoilMM:   waterNeedDryMM,
			WaterNeedDrySoilDays: waterNeedDryDays,
		}

		plantResponses = append(plantResponses, response)
	}

	respondWithJSON(http.StatusOK, plantResponses, w, cfg.sl)
	cfg.sl.Debug("User successfully listed all available plants", "user id", requestUserID)
}
