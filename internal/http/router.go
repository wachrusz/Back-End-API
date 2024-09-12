package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/wachrusz/Back-End-API/docs"
	"github.com/wachrusz/Back-End-API/internal/auth"
	"github.com/wachrusz/Back-End-API/internal/history"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	"github.com/wachrusz/Back-End-API/internal/profile"
	"net/http"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		auth.RegisterHandlers(r)
		profile.RegisterHandlers(r)
		history.RegisterHandlers(r)

		v1.RegisterHandler(r)
	})
	return r
}

func InitRouters() (chi.Router, chi.Router, error) {
	mainRouter := NewRouter()
	docRouter := docs.NewRouter()

	// Группа для изображений профиля
	mainRouter.Route("/v1/profile/image", func(r chi.Router) {
		r.Get("/get/{id}", profile.GetAvatarHandler) // Маршрут для получения изображения профиля
	})

	// Группа для эмодзи
	//mainRouter.Route("/api/emojis", func(r chi.Router) {
	//	r.Get("/get/{id}", profile.GetIconHandler) // Маршрут для получения иконки
	//})

	fmt.Println("Available routes:")
	_ = chi.Walk(mainRouter, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	return mainRouter, docRouter, nil
}
