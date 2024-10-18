package app

import (
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/config"
	api "github.com/wachrusz/Back-End-API/internal/http"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service"
	"github.com/wachrusz/Back-End-API/pkg/logger"
)

func Run(cfg *config.Config) error {
	//log.Fatalf("%s", cfg.GetDBURL())
	db, err := mydb.Init(cfg.GetDBURL())
	if err != nil {
		return err
	}
	defer db.Close()

	mydb.SetDB(db) // TODO: Избавиться от этой хуйни окончательно!

	deps := service.Dependencies{
		Repo: db,
	}

	services, err := service.NewServices(deps)
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	handler := v1.NewHandler(services)

	router, docRouter, errR := api.InitRouters(handler)
	services.Users.InitActiveUsers()

	if errR != nil {
		logger.ErrorLogger.Fatal(errR)
	}

	http.Handle("/", handler.ContentTypeMiddleware(router))
	http.Handle("/swagger/", docRouter)
	http.Handle("/docs/", docRouter)

	go services.Currency.ScheduleCurrencyUpdates()

	//changed tls hosting now everything works
	err = http.ListenAndServeTLS(":8080", cfg.CrtPath, cfg.KeyPath, nil)
	if err != nil {
		return err
	}

	return nil
}

// @Summary Get Swagger JSON
// @Description Get the Swagger JSON file.
// @Tags Swagger
// @Produce json
// @Success 200 {string} string "Swagger JSON retrieved successfully"
// @Router /swagger/json [get]
func GetSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	// Empty function, just for Swagger documentation
}

// @Summary Get Swagger UI
// @Description Get the Swagger UI.
// @Tags Swagger
// @Produce html
// @Success 200 {string} string "Swagger UI retrieved successfully"
// @Router /swagger/ui [get]
func GetSwaggerUI(w http.ResponseWriter, r *http.Request) {
	// Empty function, just for Swagger documentation
}

// @Summary Get Swagger JSON
// @Description Get the Swagger JSON file.
// @Tags Swagger
// @Produce json
// @Success 200 {string} string "Swagger JSON retrieved successfully"
// @Router /docs/swagger.json [get]
func GetSwaggerJSONFile(w http.ResponseWriter, r *http.Request) {
	// Empty function, just for Swagger documentation
}

// @Summary Get Swagger UI
// @Description Get the Swagger UI.
// @Tags Swagger
// @Produce html
// @Success 200 {string} string "Swagger UI retrieved successfully"
// @Router /swagger/index.html [get]
func GetSwaggerUIFile(w http.ResponseWriter, r *http.Request) {
	// Empty function, just for Swagger documentation
}
