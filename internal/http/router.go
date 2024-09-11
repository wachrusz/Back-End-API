package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/wachrusz/Back-End-API/internal/auth"
	"github.com/wachrusz/Back-End-API/internal/history"
	v1 "github.com/wachrusz/Back-End-API/internal/http/v1"
	"github.com/wachrusz/Back-End-API/internal/profile"
	"net/http"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	auth.RegisterHandlers(r)
	profile.RegisterHandlers(r)
	history.RegisterHandlers(r)

	v1.RegisterHandler(r)

	return r
}
