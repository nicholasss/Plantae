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
	mux.Handle("POST /api/v1/super-admin/reset-plant-types", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.resetPlantTypesHandler))))
	mux.Handle("POST /api/v1/super-admin/reset-light", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.resetLightNeedsHandler))))
	mux.Handle("POST /api/v1/super-admin/reset-water", cfg.logMW(cfg.authSuperAdminMW(http.HandlerFunc(cfg.resetWaterNeedsHandler))))

	// === admin endpoints ===

	// admin plant species endpoints
	mux.Handle("GET /api/v1/admin/plant-species", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantSpeciesViewHandler))))
	mux.Handle("POST /api/v1/admin/plant-species", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantSpeciesCreateHandler))))
	mux.Handle("PUT /api/v1/admin/plant-species/{plantSpeciesID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminReplacePlantSpeciesInfoHandler))))
	mux.Handle("DELETE /api/v1/admin/plant-species/{plantSpeciesID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminDeletePlantSpeciesHandler))))

	// admin plant names endpoints
	mux.Handle("POST /api/v1/admin/plant-names", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantNamesCreateHandler))))
	mux.Handle("GET /api/v1/admin/plant-names", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantNamesViewHandler))))
	// mux.Handle("PUT /api/v1/admin/plant-names")
	mux.Handle("DELETE /api/v1/admin/plant-names/{plantNameID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantNamesDeleteHandler))))

	// admin plant type endpoints
	mux.Handle("POST /api/v1/admin/plant-types", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantTypesCreateHandler))))
	mux.Handle("GET /api/v1/admin/plant-types", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantTypesViewHandler))))
	mux.Handle("PUT /api/v1/admin/plant-types/{plantTypeID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantTypesUpdateHandler))))
	mux.Handle("DELETE /api/v1/admin/plant-types/{plantTypeID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminPlantTypeDeleteHandler))))

	// admin set/unset plant species to plant type
	// set plant species to plant type
	mux.Handle("POST /api/v1/admin/plant-types/link/{plantTypeID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminSetPlantAsTypeHandler))))
	// unset plant species to lighting need
	mux.Handle("DELETE /api/v1/admin/plant-types/link/{plantTypeID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminUnsetPlantAsTypeHandler))))

	// admin lighting needs endpoints
	mux.Handle("POST /api/v1/admin/light", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminLightCreateHandler))))
	mux.Handle("GET /api/v1/admin/light", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminLightViewHandler))))
	mux.Handle("PUT /api/v1/admin/light/{lightID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminLightUpdateHandler))))
	mux.Handle("DELETE /api/v1/admin/light/{lightID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminLightDeleteHandler))))

	// admin set/unset plant species to lighting need
	// set plant species to lighting need
	mux.Handle("POST /api/v1/admin/light/link/{lightID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminSetPlantAsLightNeedHandler))))
	// unset plant species to lighting need
	mux.Handle("DELETE /api/v1/admin/light/link/{lightID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminUnsetPlantAsLightNeedHandler))))

	// admin watering needs endpoints
	mux.Handle("POST /api/v1/admin/water", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminWaterCreateHandler))))
	// GET /admin/water/{water id}
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
