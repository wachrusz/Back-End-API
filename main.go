// @title Cash Advisor API
// @version 1.0
// @description Backend API for managing user profiles, authentication, analytics, and more.
// @host localhost:8080
// @BasePath /
// @schemes http
// @produces json
// @consumes json
// @license MIT
// @contact.email lstwrd@yandex.com
// @contact.name Mikhail Vakhrushin
// @contact.url
// @BasePath /v1
// @SecurityDefinitions.apiKey headerKey
// @Security api_key
// @in header
// @name Authorization
// @Server http://localhost:8080
package main

import (
	auth "backEndAPI/_auth"
	handlers "backEndAPI/_handlers"
	history "backEndAPI/_history"
	mydb "backEndAPI/_mydatabase"
	profile "backEndAPI/_profile"
	user "backEndAPI/_user"
	"os"

	//"encoding/json"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type UserProfile struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

var (
	databaseURL string = "postgres://postgres:123@localhost:5432/backendapi?sslmode=disable"
	db          *mydb.Database
)

func main() {

	db, err := mydb.Init(databaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mydb.SetDB(db)

	router := mux.NewRouter().PathPrefix("/v1").Subrouter()
	docRouter := mux.NewRouter()

	auth.SetAPIKey()

	registerHandlers(router)

	docRouter.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	docRouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("/docs/swagger.json")))
	docRouter.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))

	auth.InitActiveUsers()

	http.Handle("/", router)
	http.Handle("/swagger/", docRouter)
	http.Handle("/docs/", docRouter)
	http.ListenAndServe(":8080", nil)

	log.Fatal(http.ListenAndServe(":8080", router))
	os.Exit(0)
}

func registerHandlers(router *mux.Router) {
	auth.RegisterHandlers(router)
	profile.RegisterHandlers(router)
	history.RegisterHandlers(router)
	router.HandleFunc("/auth/register", RegisterUser).Methods("POST")
	router.HandleFunc("/analytics/income", auth.AuthMiddleware(handlers.CreateIncomeHandler)).Methods("POST")
	router.HandleFunc("/analytics/expence", auth.AuthMiddleware(handlers.CreateExpenseHandler)).Methods("POST")
	router.HandleFunc("/analytics/wealth_fund", auth.AuthMiddleware(handlers.CreateWealthFundHandler)).Methods("POST")
	router.HandleFunc("/tracker/goal", auth.AuthMiddleware(handlers.CreateGoalHandler)).Methods("POST")
	router.HandleFunc("/app/category/expense", auth.AuthMiddleware(handlers.CreateExpenseCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/category/income", auth.AuthMiddleware(handlers.CreateIncomeCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/category/investment", auth.AuthMiddleware(handlers.CreateInvestmentCategoryHandler)).Methods("POST")
	router.HandleFunc("/settings/subscription", auth.AuthMiddleware(handlers.CreateSubscriptionHandler)).Methods("POST")
}

// @Summary Register user
// @Description Register a new user.
// @Tags User
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Param name query string true "Name"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error registering user"
// @Router /auth/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	name := r.FormValue("name")

	err := user.RegisterUser(username, name, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Succesfully registred"))
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
