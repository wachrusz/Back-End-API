package initialisation

import (
	auth "backEndAPI/_auth"
	handlers "backEndAPI/_handlers"
	history "backEndAPI/_history"
	profile "backEndAPI/_profile"
	report "backEndAPI/_report"

	//"encoding/json"

	"github.com/gorilla/mux"
)

func registerHandlers(router *mux.Router) {
	auth.RegisterHandlers(router)
	profile.RegisterHandlers(router)
	history.RegisterHandlers(router)

	router.HandleFunc("/app/category/expense", auth.AuthMiddleware(handlers.CreateExpenseCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/category/income", auth.AuthMiddleware(handlers.CreateIncomeCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/category/investment", auth.AuthMiddleware(handlers.CreateInvestmentCategoryHandler)).Methods("POST")
	router.HandleFunc("/app/connected-accounts/add", auth.AuthMiddleware(handlers.AddConnectedAccountHandler)).Methods("POST")
	router.HandleFunc("/app/connected-accounts/delete", auth.AuthMiddleware(handlers.DeleteConnectedAccountHandler)).Methods("DELETE")
	router.HandleFunc("/app/report", auth.AuthMiddleware(report.ExportHandler)).Methods("GET")

	router.HandleFunc("/analytics/income", auth.AuthMiddleware(handlers.CreateIncomeHandler)).Methods("POST")
	router.HandleFunc("/analytics/expense", auth.AuthMiddleware(handlers.CreateExpenseHandler)).Methods("POST")
	router.HandleFunc("/analytics/wealth_fund", auth.AuthMiddleware(handlers.CreateWealthFundHandler)).Methods("POST")

	router.HandleFunc("/tracker/goal", auth.AuthMiddleware(handlers.CreateGoalHandler)).Methods("POST")

	router.HandleFunc("/settings/subscription", auth.AuthMiddleware(handlers.CreateSubscriptionHandler)).Methods("POST")

	router.HandleFunc("/support/request", auth.AuthMiddleware(handlers.SendSupportRequestHandler)).Methods("POST")
}
