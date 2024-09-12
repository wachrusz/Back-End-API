package main

import (
	"github.com/wachrusz/Back-End-API/internal/currency"
	initialisation "github.com/wachrusz/Back-End-API/internal/http"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	"github.com/wachrusz/Back-End-API/pkg/logger"
	mydb "github.com/wachrusz/Back-End-API/pkg/mydatabase"
	secret "github.com/wachrusz/Back-End-API/secret"
	"log"
	"os"

	//"encoding/json"

	"net/http"
)

func main() {
	db, err := mydb.Init(secret.Secret.DBURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mydb.SetDB(db)

	router, docRouter, errR := initialisation.InitRouters()

	if errR != nil {
		logger.ErrorLogger.Fatal(errR)
	}

	http.Handle("/", v1.ContentTypeMiddleware(router))
	http.Handle("/swagger/", docRouter)
	http.Handle("/docs/", docRouter)

	go currency.ScheduleCurrencyUpdates()

	//changed tls hosting now everything works
	err = http.ListenAndServeTLS(":8080", secret.Secret.CrtPath, secret.Secret.KeyPath, nil)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
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
