package main

import (
	"log"
	"net/http"

	"github.com/nicholasss/plantae/internal/auth"
)

// === Middleware Functions ===

// auth super admin middleware
func (cfg *apiConfig) authSuperAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		// log.Print("Getting SuperAdmin authentication token...")
		requestToken, err := auth.GetAuthKeysValue(r.Header, "SuperAdminToken")
		if err != nil {
			log.Printf("Error with SuperAdminToken: %q", err)
			respondWithError(err, http.StatusBadRequest, w)
			return
		}

		// authenticate request
		// log.Print("Checking SuperAdmin token for authentication...")
		if ok := auth.ValidateSuperAdmin(cfg.superAdminToken, requestToken); !ok {
			respondWithError(err, http.StatusForbidden, w)
			return
		}

		log.Print("Authenticated Super Admin successfully.")
		next.ServeHTTP(w, r)
	})
}

// auth normal admin middleware
func (cfg *apiConfig) authNormalAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUserID, err := cfg.getUserIDFromToken(r)
		if err != nil {
			log.Printf("Could not authorize user in request due to: %q", err)
			respondWithError(err, http.StatusBadRequest, w)
			return
		}

		userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), requestUserID)
		if err != nil {
			log.Printf("Could not find user record for user id %q due to %q", requestUserID, err)
			respondWithError(err, http.StatusInternalServerError, w)
			return
		}

		if !userRecord.IsAdmin {
			log.Printf("Non-Admin %q [ID %q] is performing requests to admin endpoints.", userRecord.Email, requestUserID)
			respondWithError(err, http.StatusUnauthorized, w)
			return
		}

		log.Print("Authenticated normal admin successfully.")
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
