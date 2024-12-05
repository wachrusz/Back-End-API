package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/wachrusz/Back-End-API/docs"
	"github.com/wachrusz/Back-End-API/internal/http/obhttp"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	"go.uber.org/zap"
)

func newRouter(h *v1.MyHandler, obh *obhttp.MyHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(h.ContentTypeMiddleware, h.RateLimitMiddleware)

	r.Route("/v1", func(r chi.Router) {
		h.RegisterHandler(r)
		h.RegisterUserHandlers(r)
		h.RegisterProfileHandlers(r)
	})

	//Open Banking Group
	r.Route("/open-banking", func(r chi.Router) {
		obh.RegisterOBHandlers(r)
	})
	return r
}

func InitRouters(h *v1.MyHandler, obh *obhttp.MyHandler, l *zap.Logger) (chi.Router, chi.Router, error) {
	mainRouter := newRouter(h, obh)
	docRouter := docs.NewRouter()

	// Группа для изображений профиля
	mainRouter.Route("/v1/profile/image", func(r chi.Router) {
		r.Get("/get/{id}", h.GetAvatarHandler) // Маршрут для получения изображения профиля
	})

	// Группа для эмодзи
	//mainRouter.Route("/api/emojis", func(r chi.Router) {
	//	r.Get("/get/{id}", profile.GetIconHandler) // Маршрут для получения иконки
	//})

	//l.Debug("Available routes:")
	//
	//for _, r := range [2]chi.Router{mainRouter, docRouter} {
	//	_ = chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	//		l.Debug("Route found", zap.String("method", method), zap.String("route", route))
	//		return nil
	//	})
	//}

	return mainRouter, docRouter, nil
}
