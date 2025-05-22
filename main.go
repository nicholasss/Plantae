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
	mux.Handle("GET /health", cfg.logMW(http.HandlerFunc(healthHandler)))

	// user endpoints
	mux.Handle("POST /api/v1/create-user", cfg.logMW(http.HandlerFunc(cfg.createUserHandler)))

	// admin endpoints
	mux.Handle("POST /api/v1/promote-user", cfg.logMW(cfg.authenticateAdminMiddleware(http.HandlerFunc(cfg.promoteUserToAdminHandler))))
	mux.Handle("POST /api/v1/demote-user", cfg.logMW(cfg.authenticateAdminMiddleware(http.HandlerFunc(cfg.demoteUserToAdminHandler))))

	log.Printf("Server is now online at http://%s%s.\n", cfg.localAddr, cfg.port)
	log.Fatal(http.ListenAndServe(cfg.port, mux))
}
