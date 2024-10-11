package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/wachrusz/Back-End-API/docs"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	"net/http"
)

func newRouter(h *v1.MyHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		h.RegisterHandler(r)
		h.RegisterUserHandlers(r)
		//history.RegisterHandlers(r)
		h.RegisterProfileHandlers(r)
	})
	return r
}

func InitRouters(h *v1.MyHandler) (chi.Router, chi.Router, error) {
	mainRouter := newRouter(h)
	docRouter := docs.NewRouter()

	// Группа для изображений профиля
	mainRouter.Route("/v1/profile/image", func(r chi.Router) {
		r.Get("/get/{id}", h.GetAvatarHandler) // Маршрут для получения изображения профиля
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
