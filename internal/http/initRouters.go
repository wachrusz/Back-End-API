package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/wachrusz/Back-End-API/docs"
	"github.com/wachrusz/Back-End-API/internal/auth/service"
	"github.com/wachrusz/Back-End-API/internal/currency"
	"github.com/wachrusz/Back-End-API/internal/profile"
	"log"

	"net/http"
)

func InitRouters() (chi.Router, chi.Router, error) {
	mainRouter := NewRouter()
	docRouter := docs.NewRouter()

	// Группа для изображений профиля
	mainRouter.Route("/profile/image", func(r chi.Router) {
		r.Get("/get/{id}", profile.GetAvatarHandler) // Маршрут для получения изображения профиля
	})

	// Группа для эмодзи
	mainRouter.Route("/api/emojis", func(r chi.Router) {
		r.Get("/get/{id}", profile.GetIconHandler) // Маршрут для получения иконки
	})

	err := currency.InitCurrentCurrencyData()
	if err != nil {
		return nil, nil, err
	}

	//auth.SetAPIKey()

	fmt.Println("Available routes:")
	_ = chi.Walk(mainRouter, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

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
