package obhttp

import "github.com/go-chi/chi/v5"

func (h *MyHandler) RegisterOBHandlers(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/get-token", h.GetTokenHandler)
		})
		r.Route("/consents", func(r chi.Router) {
			r.Get("/get", h.GetTokenHandler)
			r.Post("/create", h.GetTokenHandler)
			r.Delete("/delete", h.GetTokenHandler)
		})
		r.Route("/accounts", func(r chi.Router) {
			r.Get("/get", h.GetAccountsHandler)
		})
	})
}
