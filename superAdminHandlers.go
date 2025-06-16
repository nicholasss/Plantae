package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// === Plant Types Management Handlers ===

// reset plant types table
func (cfg *apiConfig) resetPlantTypesHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		log.Printf("Unable to reset plant_types table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	// drop records from plant types table
	err := cfg.db.ResetPlantTypesTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset plant_types table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset plant_types table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

// === Light Needs Management Handlers ===

// reset light needs table
func (cfg *apiConfig) resetLightNeedsHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		log.Printf("Unable to reset light_needs table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	// drop records from db
	err := cfg.db.ResetLightNeedsTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset light_needs table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset light_needs table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

// === Water Needs Management Handlers ===

// reset water needs table
func (cfg *apiConfig) resetWaterNeedsHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		log.Printf("Unable to reset water_needs table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	// drop records from db
	err := cfg.db.ResetWateringNeedsTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset water_needs table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset water_needs table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

// === Plant Species Management Handlers ===

// resets plant species table
func (cfg *apiConfig) resetPlantSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		log.Printf("Unable to reset plant_species table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	// drop records from db
	err := cfg.db.ResetPlantSpeciesTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset plant_species table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset plant_species table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

// === Plant Names Management Handlers ===

// resets plant names table
func (cfg *apiConfig) resetPlantNamesHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	// ensure development platform
	if platformProduction(cfg) {
		log.Printf("Unable to reset user table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	err := cfg.db.ResetPlantNamesTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset user table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset plant_names table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

// === User/Admin management Handlers ===

// resets user table
func (cfg *apiConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	// ensure development platform
	if platformProduction(cfg) {
		log.Printf("Unable to reset user table due to platform: %q", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w)
		return
	}

	// drop records from db
	err := cfg.db.ResetUsersTable(r.Context())
	if err != nil {
		log.Printf("Unable to reset user table due to error: %q", err)
		respondWithError(nil, http.StatusInternalServerError, w)
		return
	}

	// reset successfully
	log.Print("Reset users table successfully.")
	w.WriteHeader(http.StatusNoContent)
}

// promotes user to admin
func (cfg *apiConfig) promoteUserToAdminHandler(w http.ResponseWriter, r *http.Request) {
	var adminStatusRequest AdminStatusRequest
	err := json.NewDecoder(r.Body).Decode(&adminStatusRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}
	defer r.Body.Close()

	// validate that id is a users id
	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check that user is not already admin
	if userRecord.IsAdmin {
		respondWithError(fmt.Errorf("user is already admin"), http.StatusBadRequest, w)
		return
	}

	// make user admin
	err = cfg.db.PromoteUserToAdminByID(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// successful
	w.WriteHeader(http.StatusNoContent)
}

// demotes user from admin
func (cfg *apiConfig) demoteUserToAdminHandler(w http.ResponseWriter, r *http.Request) {
	var adminStatusRequest AdminStatusRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&adminStatusRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// validate that id is a users id
	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w)
		return
	}

	// check that user is not demoted was never promoted
	if !userRecord.IsAdmin {
		respondWithError(fmt.Errorf("user is already not-admin"), http.StatusBadRequest, w)
		return
	}

	// demote user
	err = cfg.db.DemoteUserFromAdminByID(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w)
		return
	}

	// successful
	w.WriteHeader(http.StatusNoContent)
}
