// @title Cash Advisor API
// @version 1.0
// @description Backend API for managing user profiles, authentication, analytics, and more.
// @host localhost:8080
// @BasePath /v1
// @schemes https
// @produces json
// @consumes json
// @license MIT
// @contact.email lstwrd@yandex.com
// @contact.name Mikhail Vakhrushin
// @contact.url
// @Security JWT
// @securityDefinitions.JWT.type apiKey
// @securityDefinitions.JWT.name Authorization
// @securityDefinitions.JWT.in header
// @Server https://localhost:8080

package main

import (
	initialisation "main/initialisation"
	logger "main/packages/_logger"
	mydb "main/packages/_mydatabase"
	secret "main/secret"
	"main/server"

	"log"
	"os"

	//"encoding/json"

	"net/http"
)

var (
	db *mydb.Database
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

	http.Handle("/", server.ContentTypeMiddleware(router))
	http.Handle("/swagger/", docRouter)
	http.Handle("/docs/", docRouter)

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
