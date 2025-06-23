package main

import (
	"net/http"

	"github.com/nicholasss/plantae/internal/auth"
)

// === Middleware Functions ===

// auth super admin middleware
func (cfg *apiConfig) authSuperAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestToken, err := auth.GetAuthKeysValue(r.Header, "SuperAdminToken")
		if err != nil {
			cfg.sl.Debug("Unable to get superadmin token from headers", "error", err)
			respondWithError(err, http.StatusBadRequest, w)
			return
		}

		if ok := auth.ValidateSuperAdmin(cfg.superAdminToken, requestToken); !ok {
			cfg.sl.Debug("Unable to validate superadmin token in request")
			respondWithError(err, http.StatusForbidden, w)
			return
		}

		cfg.sl.Debug("Authenticated Super Admin successfully")
		next.ServeHTTP(w, r)
	})
}

// auth normal admin middleware
func (cfg *apiConfig) authNormalAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUserID, err := cfg.getUserIDFromToken(r)
		if err != nil {
			cfg.sl.Debug("Could not authorize user in request", "error", err)
			respondWithError(err, http.StatusBadRequest, w)
			return
		}

		userRecord, err := cfg.db.GetUserByIDWithoutPassword(r.Context(), requestUserID)
		if err != nil {
			cfg.sl.Debug("Could not find user record", "user id", requestUserID, "error", err)
			respondWithError(err, http.StatusInternalServerError, w)
			return
		}

		if !userRecord.IsAdmin {
			cfg.sl.Debug("Non-Admin is performing requests to admin endpoints", "email", userRecord.Email, "id", requestUserID)
			respondWithError(err, http.StatusUnauthorized, w)
			return
		}

		cfg.sl.Debug("Authenticated normal admin successfully")
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.sl.Debug("Incoming request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
