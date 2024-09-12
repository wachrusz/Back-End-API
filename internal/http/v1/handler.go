package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/wachrusz/Back-End-API/internal/auth"
	"github.com/wachrusz/Back-End-API/internal/report"
)

//type MyHandler interface {
//	GetHandler() *http.Handler
//}

func RegisterHandler(r chi.Router) {
	r.Route("/app", func(r chi.Router) {
		r.Route("/category", func(r chi.Router) {
			r.Post("/expense", auth.AuthMiddleware(CreateExpenseCategoryHandler))
			r.Post("/income", auth.AuthMiddleware(CreateIncomeCategoryHandler))
			r.Post("/investment", auth.AuthMiddleware(CreateInvestmentCategoryHandler))
		})

		r.Route("/connected-accounts", func(r chi.Router) {
			r.Post("/add", auth.AuthMiddleware(AddConnectedAccountHandler))
			r.Delete("/delete", auth.AuthMiddleware(DeleteConnectedAccountHandler))
		})

		r.Get("/report", auth.AuthMiddleware(report.ExportHandler))
	})

	r.Route("/analytics", func(r chi.Router) {
		r.Post("/income", auth.AuthMiddleware(CreateIncomeHandler))
		r.Post("/expense", auth.AuthMiddleware(CreateExpenseHandler))
		r.Post("/wealth_fund", auth.AuthMiddleware(CreateWealthFundHandler))
	})

	r.Post("/tracker/goal", auth.AuthMiddleware(CreateGoalHandler))

	r.Post("/settings/subscription", auth.AuthMiddleware(CreateSubscriptionHandler))

	r.Post("/support/request", auth.AuthMiddleware(SendSupportRequestHandler))
}
