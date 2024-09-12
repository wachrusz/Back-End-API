package docs

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	// Маршрут для статической документации
	r.Route("/docs", func(docRouter chi.Router) {
		docRouter.Get("/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))).ServeHTTP)
	})

	// Маршрут для Swagger UI
	r.Route("/swagger", func(swaggerRouter chi.Router) {
		swaggerRouter.Get("/*", httpSwagger.Handler(
			httpSwagger.URL("/docs/swagger.json"),
		).ServeHTTP)
	})
	return r
}
