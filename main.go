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

	// === super-admin endpoints ===

	mux.Handle("POST /api/v1/super-admin/promote-user", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.promoteUserToAdminHandler))))
	mux.Handle("POST /api/v1/super-admin/demote-user", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.demoteUserToAdminHandler))))

	// reset endpoints utilized for development & testing
	// requires super-admin token & for platform to be not production.

	// suepr-admin user reset endpoints
	mux.Handle("POST /api/v1/super-admin/reset-users", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.resetUsersHandler))))

	// super-admin plant reset endpoints
	mux.Handle("POST /api/v1/super-admin/reset-plant-species", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.resetPlantSpeciesHandler))))
	mux.Handle("POST /api/v1/super-admin/reset-plant-names", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.resetPlantNamesHandler))))
	// POST /super-admin/reset-plant-types
	// POST /super-admin/reset-light
	// POST /super-admin/reset-water

	// === admin endpoints ===

	// admin plant species endpoints
	mux.Handle("GET /api/v1/admin/plant-species", cfg.logMW(http.HandlerFunc(cfg.adminPlantSpeciesViewHandler)))
	mux.Handle("POST /api/v1/admin/plant-species", cfg.logMW(http.HandlerFunc(cfg.adminPlantSpeciesCreateHandler)))
	mux.Handle("PUT /api/v1/admin/plant-species/{plantSpeciesID}", cfg.logMW(http.HandlerFunc(cfg.adminReplacePlantSpeciesInfoHandler)))
	mux.Handle("DELETE /api/v1/admin/plant-species/{plantSpeciesID}", cfg.logMW(http.HandlerFunc(cfg.adminDeletePlantSpeciesHandler)))

	// admin plant names endpoints
	mux.Handle("POST /api/v1/admin/plant-names", cfg.logMW(http.HandlerFunc(cfg.adminPlantNamesCreateHandler)))
	mux.Handle("GET /api/v1/admin/plant-names", cfg.logMW(http.HandlerFunc(cfg.adminPlantNamesViewHandler)))
	mux.Handle("DELETE /api/v1/admin/plant-names/{plantNameID}", cfg.logMW(http.HandlerFunc(cfg.adminPlantNamesDeleteHandler)))

	// admin plant type endpoints
	// POST /admin/plant-type
	// GET /admin/plant-type
	// PUT /admin/plant-type/{plant type id}
	// DELETE /admin/plant-type/{plant type id}

	// admin set/unset plant species to plant type
	// set plant species to plant type
	// POST /admin/plant-type/{plant type id} ? plant species id = uuid
	// unset plant species to lighting need
	// DELETE /admin/plant-type ? plant species id = uuid

	// admin lighting needs endpoints
	// POST /admin/light
	// GET /admin/light
	// PUT /admin/light/{light id}
	// DELETE /admin/light/{light id}

	// admin set/unset plant species to lighting need
	// set plant species to lighting need
	// -- POST /admin/light/{light id} ? plant species id = uuid
	// unset plant species to lighting need
	// -- DELETE /admin/light ? plant species id = uuid

	// admin watering needs endpoints
	// POST /admin/water
	// GET /admin/water
	// PUT /admin/water/{water id}
	// DELETE /admin/water/{water id}

	// admin set/unset plant species to watering need
	// set plant species to watering need
	// -- POST /admin/water/{water id} ? plant species id = uuid
	// unset plant species to watering need
	// -- DELETE /admin/water ? plant species id = uuid

	// === user endpoints ===

	// user auth endpoints
	mux.Handle("POST /api/v1/auth/register", cfg.logMW(http.HandlerFunc(cfg.createUserHandler)))
	mux.Handle("POST /api/v1/auth/login", cfg.logMW(http.HandlerFunc(cfg.loginUserHandler)))
	mux.Handle("POST /api/v1/auth/refresh", cfg.logMW(http.HandlerFunc(cfg.refreshUserHandler)))
	mux.Handle("POST /api/v1/auth/revoke", cfg.logMW(http.HandlerFunc(cfg.revokeUserHandler)))

	// === user data endpoints
	// mux.Handle("GET /api/v1/my/plants", cfg.logMW(http.HandlerFunc(cfg.usersPlantsListHandler)))

	// additional tables/columns to add?
	// pot type, pot date, aquisition date

	// POST /my/plants
	// GET /my/plants
	// PUT /my/plants/{plant id}
	// DELETE /my/plants/{plant id}
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
