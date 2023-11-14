package main

import (
	auth "backEndAPI/_auth"
	categories "backEndAPI/_categories"
	history "backEndAPI/_history"
	mydb "backEndAPI/_mydatabase"
	profile "backEndAPI/_profile"
	user "backEndAPI/_user"

	//"encoding/json"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Простой тип для представления профиля пользователя
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

	auth.SetDB(db)
	user.SetDB(db)

	router := mux.NewRouter()

	auth.SetAPIKey("123")

	registerHandlers(router)

	// Запуск сервера на порту 8080
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func registerHandlers(router *mux.Router) {
	auth.RegisterHandlers(router)
	profile.RegisterHandlers(router)
	history.RegisterHandlers(router)
	categories.RegisterHandlers(router)
	router.HandleFunc("/user/register", RegisterUser).Methods("POST")
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Получаем данные из запроса
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Регистрируем пользователя
	err := user.RegisterUser(username, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Succesfully registred"))
}
