package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// === Handler Functions ===

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// === Main Function ===

func main() {
	log.Printf("Staring server\n")

	cfg, err := loadApiConfig()
	if err != nil {
		log.Fatalf("Issue loading config: %q", err)
	}

	mux := http.NewServeMux()

	// health endpoint
	mux.Handle("GET /api/v1/health", cfg.logMW(http.HandlerFunc(healthHandler)))

	// super-admin endpoints
	mux.Handle("POST /api/v1/super-admin/reset-users", cfg.logMW(cfg.authenticateAdminMiddleware(http.HandlerFunc(cfg.resetUsersHandler))))
	mux.Handle("POST /api/v1/super-admin/promote-user", cfg.logMW(cfg.authenticateAdminMiddleware(http.HandlerFunc(cfg.promoteUserToAdminHandler))))
	mux.Handle("POST /api/v1/super-admin/demote-user", cfg.logMW(cfg.authenticateAdminMiddleware(http.HandlerFunc(cfg.demoteUserToAdminHandler))))

	// user auth endpoints
	mux.Handle("POST /api/v1/auth/register", cfg.logMW(http.HandlerFunc(cfg.createUserHandler)))
	mux.Handle("POST /api/v1/auth/login", cfg.logMW(http.HandlerFunc(cfg.loginUserHandler)))
	mux.Handle("POST /api/v1/auth/refresh", cfg.logMW(http.HandlerFunc(cfg.refreshUserHandler)))
	mux.Handle("POST /api/v1/auth/revoke", cfg.logMW(http.HandlerFunc(cfg.revokeUserHandler)))

	// === user data endpoints
	//
	// /my/plants
	// /my/plants/{plant id}
	//
	// /my/rooms
	// /my/rooms/{room id}
	//
	// === data endpoints
	//
	// /biomes
	// /biomes/{biome id}
	//
	// /plants
	// /plants/{plant id}

	log.Printf("Server is now online at http://%s%s.\n", cfg.localAddr, cfg.port)
	log.Fatal(http.ListenAndServe(cfg.port, mux))
}
