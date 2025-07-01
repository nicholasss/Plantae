package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// === request response types ===

// AdminStatusRequest is for promotion or demotion requests of a user, from an super-admin.
type AdminStatusRequest struct {
	ID uuid.UUID `json:"id"`
}

type AdminStatusResponse struct {
	ID      uuid.UUID `json:"id"`
	IsAdmin bool      `json:"isAdmin"`
}

// === Plant Types Management Handlers ===

// reset plant types table
// POST /api/v1/super-admin/reset-plant-types
// 204 No Content is ok in context
func (cfg *apiConfig) resetPlantTypesHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		cfg.sl.Debug("Unable to reset plant_types table due to wrong platform", "platform", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w, cfg.sl)
		return
	}

	// drop records from plant types table
	err := cfg.db.ResetPlantTypesTable(r.Context())
	if err != nil {
		cfg.sl.Debug("Unable to reset plant_types table", "error", err)
		respondWithError(nil, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Reset plant_types table successfully")
	w.WriteHeader(http.StatusNoContent)
}

// === Light Needs Management Handlers ===

// reset light needs table
// POST /api/v1/super-admin/reset-light
// 204 No Content is ok in context
func (cfg *apiConfig) resetLightNeedsHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		cfg.sl.Debug("Unable to reset light_needs table due to wrong platform", "platform", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w, cfg.sl)
		return
	}

	// drop records from db
	err := cfg.db.ResetLightNeedsTable(r.Context())
	if err != nil {
		cfg.sl.Debug("Unable to reset light_needs table", "error", err)
		respondWithError(nil, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Reset light_needs table successfully")
	w.WriteHeader(http.StatusNoContent)
}

// === Water Needs Management Handlers ===

// reset water needs table
// POST /api/v1/super-admin/reset-water
// 204 No Content is ok in context
func (cfg *apiConfig) resetWaterNeedsHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		cfg.sl.Debug("Unable to reset water_needs table due to wrong platform", "platform", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w, cfg.sl)
		return
	}

	// drop records from db
	err := cfg.db.ResetWaterNeedsTable(r.Context())
	if err != nil {
		cfg.sl.Debug("Unable to reset water_needs table", "error", err)
		respondWithError(nil, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Reset water_needs table successfully")
	w.WriteHeader(http.StatusNoContent)
}

// === Plant Species Management Handlers ===

// resets plant species table
// POST /api/v1/super-admin/reset-plant-species
// 204 No Content is ok in context
func (cfg *apiConfig) resetPlantSpeciesHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	if platformProduction(cfg) {
		cfg.sl.Debug("Unable to reset plant_species table due to wrong platform", "platform", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w, cfg.sl)
		return
	}

	// drop records from db
	err := cfg.db.ResetPlantSpeciesTable(r.Context())
	if err != nil {
		cfg.sl.Debug("Unable to reset plant_species table", "error", err)
		respondWithError(nil, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Reset plant_species table successfully")
	w.WriteHeader(http.StatusNoContent)
}

// === Plant Names Management Handlers ===

// resets plant names table
// POST /api/v1/super-admin/reset-plant-names
// 204 No Content is ok in context
func (cfg *apiConfig) resetPlantNamesHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	// ensure development platform
	if platformProduction(cfg) {
		cfg.sl.Debug("Unable to reset plant_names table due to wrong platform", "platform", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w, cfg.sl)
		return
	}

	err := cfg.db.ResetPlantNamesTable(r.Context())
	if err != nil {
		cfg.sl.Debug("Unable to reset plant_names table", "error", err)
		respondWithError(nil, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Reset plant_names table successfully")
	w.WriteHeader(http.StatusNoContent)
}

// === User/Admin management Handlers ===

// resets user table
// POST /api/v1/super-admin/reset-users
// 204 No Content is ok in context
func (cfg *apiConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// super-admin pre-authenticated before the handler is used
	// ensure development platform
	if platformProduction(cfg) {
		cfg.sl.Debug("Unable to reset users table due to wrong platform", "platform", cfg.platform)
		respondWithError(nil, http.StatusForbidden, w, cfg.sl)
		return
	}

	// drop records from db
	err := cfg.db.ResetUsersTable(r.Context())
	if err != nil {
		cfg.sl.Debug("Unable to reset users table", "error", err)
		respondWithError(nil, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Reset users table successfully")
	w.WriteHeader(http.StatusNoContent)
}

// promotes user to admin
// POST /api/v1/super-admin/promote-user
// 200 OK makes sense in context
func (cfg *apiConfig) promoteUserToAdminHandler(w http.ResponseWriter, r *http.Request) {
	var adminStatusRequest AdminStatusRequest
	err := json.NewDecoder(r.Body).Decode(&adminStatusRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}
	defer r.Body.Close()

	// validate that id is a users id
	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// check that user is not already admin
	if userRecord.IsAdmin {
		respondWithError(errors.New("user is already admin"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	err = cfg.db.PromoteUserToAdminByID(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// send response body
	adminResponse := AdminStatusResponse{
		ID:      userRecord.ID,
		IsAdmin: true,
	}
	adminData, err := json.Marshal(adminResponse)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Successfully promoted user to admin", "user id", userRecord.ID)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(adminData)
}

// demotes user from admin
// POST /api/v1/super-admin/demote-user
// 200 OK makes sense in context
func (cfg *apiConfig) demoteUserToAdminHandler(w http.ResponseWriter, r *http.Request) {
	var adminStatusRequest AdminStatusRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&adminStatusRequest)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// validate that id is a users id
	userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// check that user is not demoted was never promoted
	if !userRecord.IsAdmin {
		respondWithError(errors.New("user is already not-admin"), http.StatusBadRequest, w, cfg.sl)
		return
	}

	err = cfg.db.DemoteUserFromAdminByID(r.Context(), adminStatusRequest.ID)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	// send response body
	adminResponse := AdminStatusResponse{
		ID:      userRecord.ID,
		IsAdmin: false,
	}
	adminData, err := json.Marshal(adminResponse)
	if err != nil {
		cfg.sl.Debug("Could not marshal data", "error", err)
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	cfg.sl.Info("Successfully demoted user from admin", "user id", userRecord.ID)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(adminData)
}
