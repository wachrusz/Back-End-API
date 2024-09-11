package http

import (
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/auth/service"
	"github.com/wachrusz/Back-End-API/internal/currency"
	"github.com/wachrusz/Back-End-API/internal/profile"
	"log"
	//"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func InitRouters() (*mux.Router, *mux.Router, error) {
	mainRouter := mux.NewRouter().PathPrefix("/v1").Subrouter()
	docRouter := mux.NewRouter()
	registerHandlers(mainRouter)

	imageGroup := mainRouter.PathPrefix("/profile/image").Subrouter()
	{
		imageGroup.Methods(http.MethodGet).Path("/get/{id}").HandlerFunc(profile.GetAvatarHandler)
	}

	iconGroup := mainRouter.PathPrefix("/api/emojis").Subrouter()
	{
		iconGroup.Methods(http.MethodGet).Path("/get/{id}").HandlerFunc(profile.GetIconHandler)
	}

	err := currency.InitCurrentCurrencyData()
	if err != nil {
		return nil, nil, err
	}

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

	service.InitActiveUsers()

	baseURL := getBaseURL()

	log.Println("Base URL:", baseURL)

	return mainRouter, docRouter, nil
}

func getBaseURL() string {
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic(err)
	}

	baseURL := request.URL.String()

	return baseURL
}
