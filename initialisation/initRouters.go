package initialisation

import (
	auth "backEndAPI/_auth"
	"log"

	"fmt"

	//"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouters() (*mux.Router, *mux.Router, error) {
	mainRouter := mux.NewRouter().PathPrefix("/v1").Subrouter()
	docRouter := mux.NewRouter()
	registerHandlers(mainRouter)

	//auth.SetAPIKey()

	fmt.Println("Available routes:")
	errR := mainRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("Route:", path)
		}
		return err
	})

	if errR != nil {
		fmt.Println("Error:", errR)
		return nil, nil, errR
	}

	docRouter.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	docRouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("/docs/swagger.json")))
	docRouter.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
	))

	auth.InitActiveUsers()

	baseURL := getBaseURL()

	log.Println("Base URL:", baseURL)

	return mainRouter, docRouter, nil
}

func getBaseURL() string {
	// Используем стандартный пакет net/http для определения базовой ссылки
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}

	// Запрос URL из запроса
	baseURL := request.URL.String()

	// Получаем строковое представление базовой ссылки
	return baseURL
}
