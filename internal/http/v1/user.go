package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	email_conf "github.com/wachrusz/Back-End-API/internal/email"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/validator"
	"log"
	"net/http"
	"time"
)

// @Summary Register user
// @Description Register a new user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Param name query string true "Name"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error registering user"
// @Router /auth/register [post]
func (h *MyHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var registrationRequest user.UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&registrationRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid request payload: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// validation

	if !validator.IsValidEmail(registrationRequest.Email) {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid email: %s", registrationRequest.Email), http.StatusBadRequest)
		return
	}

	if !validator.IsValidPassword(registrationRequest.Password) {
		log.Println("invalid password")
		jsonresponse.SendErrorResponse(w, fmt.Errorf("password must be at least 7 digits long"), http.StatusBadRequest)
		return
	}

	token, err := h.s.Users.PrimaryRegistration(registrationRequest.Email, registrationRequest.Password)
	if err != nil {
		switch err {
		case user.ErrDuplicated:

		default:
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid request payload: %v", err.Error()), http.StatusBadRequest)
		}
		return
	}

	response := map[string]interface{}{
		"message":     "Confirm your email",
		"token":       token,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MyHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var resetRequest user.ResetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	if !validator.IsValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	err = h.s.Users.ResetPassword(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
	}
}

func (h *MyHandler) ChangePasswordForRecoverHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var resetRequest user.UserPasswordReset

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	password := resetRequest.Password
	if !validator.IsValidPassword(password) {
		jsonresponse.SendErrorResponse(w, errors.New("password must be at least 7 digits long: "), http.StatusBadRequest)
		return
	}

	if !validator.IsValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	resetToken := resetRequest.ResetToken
	err = h.s.Users.ChangePasswordForRecover(email, password, resetToken)
	if err != nil {
		var statusCode = 500
		switch {
		case errors.Is(err, user.ErrEmpty):
			statusCode = http.StatusBadRequest
			break
		case errors.Is(err, user.ErrEmailing):
			statusCode = http.StatusInternalServerError
			break
		case errors.Is(err, user.ErrInvalidToken):
			statusCode = http.StatusUnauthorized
			break
		}
		jsonresponse.SendErrorResponse(w, fmt.Errorf("error changing the password: %v", err), statusCode)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfuly reseted password",
		"status_code": http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Login to the system
// @Description Login to the system and get an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body auth.UserAuthenticationRequest true "UserAuthenticationRequest object"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}

	var loginRequest user.UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := loginRequest.Email
	password := loginRequest.Password

	token, err := h.s.Users.Login(email, password)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmpty) || errors.Is(err, user.ErrInvalidCreds):
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid email or password: %w", err), http.StatusUnauthorized)
		case errors.Is(err, user.ErrInternal) || errors.Is(err, user.ErrEmailing):
			jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"message":     "Confirm your email",
		"token":       token,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func (h *MyHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}
	//! Заставляет задуматься
	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	//! Сомнительно
	type token struct {
		RefreshToken string `json:"refresh_token"`
	}

	var tmpToken token

	err := json.NewDecoder(r.Body).Decode(&tmpToken)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	access, refresh, err := h.s.Users.RefreshToken(tmpToken.RefreshToken, userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":                 "Successfully refreshed tokens",
		"access_token":            access,
		"refresh_token":           refresh,
		"access_token_life_time":  time.Minute * 15,
		"refresh_token_life_time": 30 * 24 * time.Hour,
		"status_code":             http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

// @Summary Logout the user
// @Description Logs out the user, terminating the session.
// @Tags Auth
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Failure 500 {string} string "Internal Server Error"
// @Security JWT
// @Router /auth/logout [post]
func (h *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	currentDeviceID, ok := getDeviceIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	err := user.Logout(currentDeviceID, userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Logout Successful",
		"status_code": http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func getDeviceIDFromContext(ctx context.Context) (string, bool) {
	deviceID, ok := ctx.Value("device_id").(string)
	return deviceID, ok
}

func getUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

func (h *MyHandler) RegisterHandlers(router chi.Router) {
	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/login/confirm", email.ConfirmEmailLoginHandler)
		r.Post("/logout", AuthMiddleware(h.Logout))
		r.Post("/register", h.RegisterUserHandler)
		r.Post("/register/confirm-email", email.ConfirmEmailHandler)

		// Password reset routes
		r.Route("/login/reset", func(r chi.Router) {
			r.Post("/password", h.ResetPasswordHandler)
			r.Post("/password/confirm", email.ResetPasswordConfirmHandler)
			r.Put("/password/put", h.ChangePasswordForRecoverHandler)
		})

		// Token routes
		r.Post("/refresh", AuthMiddleware(h.RefreshTokenHandler))
		r.Delete("/tokens/delete", DeleteTokensHandler)
		r.Get("/tokens/amount", GetTokenPairsAmmountHandler)
	})

	// OAuth login routes
	router.Route("/auth/login", func(r chi.Router) {
		r.Get("/vk", h.s.Users.HandleVKLogin)
		r.Get("/google", h.s.Users.HandleGoogleLogin)
	})

	// Developer routes
	router.Get("/dev/confirmation-code/get", email_conf.GetConfirmationCodeTestHandler)
}
