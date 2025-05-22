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
		requestToken, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(err, http.StatusBadRequest, w)
			return
		}

		// authenticate request
		if ok := auth.ValidateSuperAdmin(cfg.superAdminToken, requestToken); !ok {
			respondWithError(err, http.StatusForbidden, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
