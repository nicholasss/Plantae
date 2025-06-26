package main

import (
	"encoding/json"
	"net/http"

	"github.com/nicholasss/plantae/internal/auth"
)

// requires access token in auth header
// creates a user_plant
func (cfg *apiConfig) usersPlantsCreateHandler(w http.ResponseWriter, r *http.Request) {
	// todo
}

// requires access token in auth header
// returns the users list of plants
func (cfg *apiConfig) usersPlantsListHandler(w http.ResponseWriter, r *http.Request) {
	accessTokenProvided, err := auth.GetBearerToken(r.Header, cfg.sl)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	requestUserID, err := auth.ValidateJWT(accessTokenProvided, cfg.JWTSecret, cfg.sl)
	if err != nil {
		respondWithError(err, http.StatusBadRequest, w, cfg.sl)
		return
	}

	// get list of plants in user_plants table
	usersPlants, err := cfg.db.GetAllUsersPlantsOrderedByUpdated(r.Context(), requestUserID)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	if len(usersPlants) <= 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	usersPlantsData, err := json.Marshal(usersPlants)
	if err != nil {
		respondWithError(err, http.StatusInternalServerError, w, cfg.sl)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(usersPlantsData)
}
