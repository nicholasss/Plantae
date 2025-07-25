package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var logo = `
$$$$$$$\  $$\                      $$\
$$  __$$\ $$ |                     $$ |
$$ |  $$ |$$ | $$$$$$\  $$$$$$$\ $$$$$$\    $$$$$$\   $$$$$$\
$$$$$$$  |$$ | \____$$\ $$  __$$\\_$$  _|   \____$$\ $$  __$$\
$$  ____/ $$ | $$$$$$$ |$$ |  $$ | $$ |     $$$$$$$ |$$$$$$$$ |
$$ |      $$ |$$  __$$ |$$ |  $$ | $$ |$$\ $$  __$$ |$$   ____|
$$ |      $$ |\$$$$$$$ |$$ |  $$ | \$$$$  |\$$$$$$$ |\$$$$$$$\
\__|      \__| \_______|\__|  \__|  \____/  \_______| \_______|`

// === Handler Functions ===

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// === Main Function ===

func main() {
	fmt.Printf("%s\n\n", logo)

	cfg, closeLogFile, err := loadAPIConfig()
	if err != nil {
		log.Fatalf("Issue loading config: %q", err)
	}
	defer closeLogFile()
	cfg.sl.Info("Starting server...")

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
	mux.Handle("GET /api/v1/admin/water", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminWaterViewHandler))))
	mux.Handle("DELETE /api/v1/admin/water/{waterID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminWaterDeleteHandler))))

	// admin set/unset plant species to watering need
	// set plant species to watering need
	mux.Handle("POST /api/v1/admin/water/link/{waterID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminSetPlantAsWaterNeedHandler))))
	// unset plant species to watering need
	mux.Handle("DELETE /api/v1/admin/water/link/{waterID}", cfg.logMW(cfg.authNormalAdminMW(http.HandlerFunc(cfg.adminUnsetPlantAsWaterNeedHandler))))

	// === user endpoints ===

	// user auth endpoints
	mux.Handle("POST /api/v1/auth/register", cfg.logMW(http.HandlerFunc(cfg.registerUserHandler)))
	mux.Handle("POST /api/v1/auth/login", cfg.logMW(http.HandlerFunc(cfg.loginHandler)))
	mux.Handle("POST /api/v1/auth/refresh", cfg.logMW(http.HandlerFunc(cfg.refreshTokenHandler)))
	mux.Handle("POST /api/v1/auth/revoke", cfg.logMW(http.HandlerFunc(cfg.revokeRefreshTokenHandler)))

	// === user data endpoints
	mux.Handle("GET /api/v1/my/plants", cfg.logMW(http.HandlerFunc(cfg.usersPlantsListHandler)))
	mux.Handle("POST /api/v1/my/plants", cfg.logMW(http.HandlerFunc(cfg.usersPlantsCreateHandler)))
	mux.Handle("PUT /api/v1/my/plants/{plantID}", cfg.logMW(http.HandlerFunc(cfg.userPlantsUpdateHandler)))
	mux.Handle("DELETE /api/v1/my/plants/{plantID}", cfg.logMW(http.HandlerFunc(cfg.userPlantsDeleteHandler)))

	// listing all plants on the server
	mux.Handle("GET /api/v1/plants", cfg.logMW(http.HandlerFunc(cfg.usersViewPlantsListHandler)))

	serverAddress := fmt.Sprintf("http://%s%s", cfg.localAddr, cfg.port)
	cfg.sl.Info("Server is now online", "address", serverAddress)
	log.Fatal(http.ListenAndServe(cfg.port, mux))
}
