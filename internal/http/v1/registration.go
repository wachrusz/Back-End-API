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

		r.Route("/accounts", func(r chi.Router) {
			r.Post("/", h.AuthMiddleware(h.AddConnectedAccountHandler))
			r.Delete("/", h.AuthMiddleware(h.DeleteConnectedAccountHandler))
			r.Put("/", h.AuthMiddleware(h.UpdateConnectedAccountHandler))
		})
	})

	r.Route("/analytics", func(r chi.Router) {
		r.Route("/income", func(r chi.Router) {
			r.Post("/", h.AuthMiddleware(h.CreateIncomeHandler))
			r.Put("/", h.AuthMiddleware(h.UpdateIncomeHandler))
			r.Delete("/", h.AuthMiddleware(h.DeleteIncomeHandler))
		})

		r.Route("/expense", func(r chi.Router) {
			r.Post("/", h.AuthMiddleware(h.CreateExpenseHandler))
			r.Put("/", h.AuthMiddleware(h.UpdateExpenseHandler))
			r.Delete("/", h.AuthMiddleware(h.DeleteExpenseHandler))
		})

		r.Route("/wealth_fund", func(r chi.Router) {
			r.Post("/", h.AuthMiddleware(h.CreateWealthFundHandler))
			r.Put("/", h.AuthMiddleware(h.UpdateWealthFundHandler))
			r.Delete("/", h.AuthMiddleware(h.DeleteWealthFundHandler))
		})
	})

	r.Route("/tracker/goal", func(r chi.Router) {
		r.Post("/", h.AuthMiddleware(h.CreateGoalHandler))
		r.Put("/", h.AuthMiddleware(h.UpdateGoalHandler))
		r.Delete("/", h.AuthMiddleware(h.DeleteGoalHandler))
	})

	r.Route("/settings/subscription", func(r chi.Router) {
		r.Post("/", h.AuthMiddleware(h.CreateSubscriptionHandler))
		r.Put("/", h.AuthMiddleware(h.UpdateSubscriptionHandler))
		r.Delete("/", h.AuthMiddleware(h.DeleteSubscriptionHandler))
	})

	r.Post("/support/request", h.AuthMiddleware(h.SendSupportRequestHandler))
}

func (h *MyHandler) RegisterUserHandlers(router chi.Router) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.LoginUserHandler)
		r.Post("/login/confirm", h.ConfirmEmailLoginHandler)
		r.Post("/logout", h.AuthMiddleware(h.LogoutUserHandler))
		r.Post("/register", h.RegisterUserHandler)
		r.Post("/register/confirm", h.ConfirmEmailRegisterHandler)

		// Password reset routes
		r.Route("/password", func(r chi.Router) {
			r.Post("/", h.ResetPasswordHandler)
			r.Post("/confirm", h.ResetPasswordConfirmHandler)
			r.Put("/", h.ChangePasswordForRecoverHandler)
		})

		// Token routes
		r.Post("/refresh", h.RefreshTokenHandler)
		r.Delete("/tokens", h.DeleteTokensHandler)
		r.Get("/tokens/amount", h.GetTokenPairsAmountHandler)
	})

	// OAuth login routes
	//router.Route("/auth/login", func(r chi.Router) {
	//	r.Get("/vk", h.s.Users.HandleVKLogin)
	//	r.Get("/google", h.s.Users.HandleGoogleLogin)
	//})

	router.Get("/metrics/rps", h.GetRequestCountHandler)

	// Developer routes
	router.Get("/dev/confirmation-code/get", h.GetConfirmationCodeTestHandler)
}

func (h *MyHandler) RegisterProfileHandlers(router chi.Router) {
	// Profile routes
	router.Route("/profile", func(r chi.Router) {
		r.Get("/", h.AuthMiddleware(h.GetProfileHandler))
		r.Get("/analytics", h.AuthMiddleware(h.GetProfileAnalyticsHandler))
		r.Get("/tracker", h.AuthMiddleware(h.GetProfileTrackerHandler))
		r.Get("/more", h.AuthMiddleware(h.GetProfileMore))
		r.Put("/name", h.AuthMiddleware(h.UpdateName))
		r.Get("/archive", h.AuthMiddleware(h.GetOperationArchive))
		//r.Put("/image/put", AuthMiddleware(user.UploadAvatarHandler))
	})

	// Emojis routes
	router.Route("/api/emojis", func(r chi.Router) {
		r.Put("/put", h.UploadIconHandler)
		r.Get("/get/list", h.GetIconsURLsHandler)
	})
}
