package main

import (
	"fmt"
	"net/http"

	"github.com/nicholasss/plantae/internal/auth"
)

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
	}

	// user is now authenticated below here
}
