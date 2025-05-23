package main

import (
	"log"
	"net/http"

	"github.com/nicholasss/plantae/internal/auth"
)

// === Middleware Functions ===

func (cfg *apiConfig) authenticateAdminMiddleware(next http.Handler) http.Handler {
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

func (cfg *apiConfig) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
