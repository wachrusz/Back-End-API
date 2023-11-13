package main

import (
	auth "backEndAPI/_auth"
	categories "backEndAPI/_categories"
	history "backEndAPI/_history"
	profile "backEndAPI/_profile"

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

var userProfile = make(map[string]UserProfile) // Имитация базы данных

func main() {
	router := mux.NewRouter()
	initializeData()

	// Регистрация обработчиков
	registerHandlers(router)

	// Запуск сервера на порту 8080
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeData() {
	// Инициализация примера данных (имитация базы данных)
	userProfile["user1"] = UserProfile{Username: "user1", Name: "John Doe"}
}

func registerHandlers(router *mux.Router) {
	auth.RegisterHandlers(router)
	profile.RegisterHandlers(router)
	history.RegisterHandlers(router)
	categories.RegisterHandlers(router)
}
