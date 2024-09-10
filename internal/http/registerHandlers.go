package http

import (
	"main/internal/_history"
	"main/internal/auth"
	handlers2 "main/internal/http/v1"
	"main/internal/profile"
	"main/internal/report"
	//"encoding/json"

	"github.com/gorilla/mux"
)

func registerHandlers(router *mux.Router) {
	auth.RegisterHandlers(router)
	profile.RegisterHandlers(router)
	history.RegisterHandlers(router)

	router.HandleFunc("/app/category/expense", auth.AuthMiddleware(handlers2.CreateExpenseCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/category/income", auth.AuthMiddleware(handlers2.CreateIncomeCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/category/investment", auth.AuthMiddleware(handlers2.CreateInvestmentCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/connected-accounts/add", auth.AuthMiddleware(handlers2.AddConnectedAccountHandler)).Methods("POST")
	router.HandleFunc("/app/connected-accounts/delete", auth.AuthMiddleware(handlers2.DeleteConnectedAccountHandler)).Methods("DELETE")
	router.HandleFunc("/app/report", auth.AuthMiddleware(report.ExportHandler)).Methods("GET")

	router.HandleFunc("/analytics/income", auth.AuthMiddleware(handlers2.CreateIncomeHandler)).Methods("POST")
	router.HandleFunc("/analytics/expense", auth.AuthMiddleware(handlers2.CreateExpenseHandler)).Methods("POST")
	router.HandleFunc("/analytics/wealth_fund", auth.AuthMiddleware(handlers2.CreateWealthFundHandler)).Methods("POST")

	router.HandleFunc("/tracker/goal", auth.AuthMiddleware(handlers2.CreateGoalHandler)).Methods("POST")

	router.HandleFunc("/settings/subscription", auth.AuthMiddleware(handlers2.CreateSubscriptionHandler)).Methods("POST")

	router.HandleFunc("/support/request", auth.AuthMiddleware(handlers2.SendSupportRequestHandler)).Methods("POST")
}
