package v1

import (
	"github.com/go-chi/chi/v5"
)

func (h *MyHandler) RegisterHandler(r chi.Router) {
	r.Route("/app", func(r chi.Router) {
		r.Route("/category", func(r chi.Router) {
			r.Post("/expense", h.AuthMiddleware(h.CreateExpenseCategoryHandler))
			r.Post("/income", h.AuthMiddleware(h.CreateIncomeCategoryHandler))
			r.Post("/investment", h.AuthMiddleware(h.CreateInvestmentCategoryHandler))
		})

		r.Route("/connected-accounts", func(r chi.Router) {
			r.Post("/add", h.AuthMiddleware(h.AddConnectedAccountHandler))
			r.Delete("/delete", h.AuthMiddleware(h.DeleteConnectedAccountHandler))
		})

		//r.Get("/report", h.AuthMiddleware(h.ExportHandler))
	})

	r.Route("/analytics", func(r chi.Router) {
		r.Post("/income", h.AuthMiddleware(h.CreateIncomeHandler))
		r.Post("/expense", h.AuthMiddleware(h.CreateExpenseHandler))
		r.Post("/wealth_fund", h.AuthMiddleware(h.CreateWealthFundHandler))
	})

	r.Post("/tracker/goal", h.AuthMiddleware(h.CreateGoalHandler))

	r.Post("/settings/subscription", h.AuthMiddleware(h.CreateSubscriptionHandler))

	r.Post("/support/request", h.AuthMiddleware(h.SendSupportRequestHandler))
}

func (h *MyHandler) RegisterUserHandlers(router chi.Router) {
	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/login/confirm", h.ConfirmEmailLoginHandler)
		r.Post("/logout", h.AuthMiddleware(h.Logout))
		r.Post("/register", h.RegisterUserHandler)
		r.Post("/register/confirm-email", h.ConfirmEmailHandler)

		// Password reset routes
		r.Route("/login/reset", func(r chi.Router) {
			r.Post("/password", h.ResetPasswordHandler)
			r.Post("/password/confirm", h.ResetPasswordConfirmHandler)
			r.Put("/password/put", h.ChangePasswordForRecoverHandler)
		})

		// Token routes
		r.Post("/refresh", h.AuthMiddleware(h.RefreshTokenHandler))
		r.Delete("/tokens/delete", h.DeleteTokensHandler)
		r.Get("/tokens/amount", h.GetTokenPairsAmountHandler)
	})

	// OAuth login routes
	//router.Route("/auth/login", func(r chi.Router) {
	//	r.Get("/vk", h.s.Users.HandleVKLogin)
	//	r.Get("/google", h.s.Users.HandleGoogleLogin)
	//})

	// Developer routes
	router.Get("/dev/confirmation-code/get", h.GetConfirmationCodeTestHandler)
}

func (h *MyHandler) RegisterProfileHandlers(router chi.Router) {
	// Profile routes
	router.Route("/profile", func(r chi.Router) {
		r.Get("/info/get", h.AuthMiddleware(h.GetProfileHandler))
		r.Get("/analytics/get", h.AuthMiddleware(h.GetProfileAnalyticsHandler))
		r.Get("/tracker/get", h.AuthMiddleware(h.GetProfileTrackerHandler))
		r.Get("/more/get", h.AuthMiddleware(h.GetProfileMore))
		r.Put("/name/put", h.AuthMiddleware(h.UpdateName))
		r.Get("/operation-archive/get", h.AuthMiddleware(h.GetOperationArchive))
		//r.Put("/image/put", AuthMiddleware(user.UploadAvatarHandler))
	})

	// Emojis routes
	router.Route("/api/emojis", func(r chi.Router) {
		r.Put("/put", h.UploadIconHandler)
		r.Get("/get/list", h.GetIconsURLsHandler)
	})
}
