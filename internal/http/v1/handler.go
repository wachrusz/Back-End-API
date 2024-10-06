package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/wachrusz/Back-End-API/internal/service"
	"net/http"
)

type MyHandler struct {
	s *service.Services
}

func NewHandler(s *service.Services) *MyHandler {
	return &MyHandler{s: s}
}

func (h *MyHandler) RegisterHandler(r chi.Router) {
	r.Route("/app", func(r chi.Router) {
		r.Route("/category", func(r chi.Router) {
			r.Post("/expense", AuthMiddleware(CreateExpenseCategoryHandler))
			r.Post("/income", AuthMiddleware(CreateIncomeCategoryHandler))
			r.Post("/investment", AuthMiddleware(CreateInvestmentCategoryHandler))
		})

		r.Route("/connected-accounts", func(r chi.Router) {
			r.Post("/add", AuthMiddleware(AddConnectedAccountHandler))
			r.Delete("/delete", AuthMiddleware(DeleteConnectedAccountHandler))
		})

		r.Get("/report", AuthMiddleware(h.ExportHandler))
	})

	r.Route("/analytics", func(r chi.Router) {
		r.Post("/income", AuthMiddleware(CreateIncomeHandler))
		r.Post("/expense", AuthMiddleware(CreateExpenseHandler))
		r.Post("/wealth_fund", AuthMiddleware(CreateWealthFundHandler))
	})

	r.Post("/tracker/goal", AuthMiddleware(CreateGoalHandler))

	r.Post("/settings/subscription", AuthMiddleware(CreateSubscriptionHandler))

	r.Post("/support/request", AuthMiddleware(SendSupportRequestHandler))
}

func (h *MyHandler) getDeviceIDFromRequest(r *http.Request) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil.String(), err
	}
	return id.String(), nil
}
