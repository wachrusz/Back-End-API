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

	router := mux.NewRouter()

	auth.SetAPIKey("123")

	registerHandlers(router)

	log.Println("Registered Routes:")
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		template, err := route.GetPathTemplate()
		if err == nil {
			log.Println(template)
		}
		return nil
	})
	auth.InitActiveUsers()

	http.Handle("/", router)
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
