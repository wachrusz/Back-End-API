package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	"github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/util"
	"github.com/wachrusz/Back-End-API/pkg/validator"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// RegisterUserHandler registers a new user in the system.
//
// @Summary Register a new user
// @Description Register a new user in the system.
// @Tags Auth
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Param name query string true "Name"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 409 {string} string "User already exists"
// @Router /auth/register [post]
func (h *MyHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Received request to register a new user.")

	// Check Content-Type header
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := fmt.Errorf("empty 'Content-Type' header")
		h.errResp(w, fmt.Errorf("invalid Content-Type, expected application/json: %v", err), http.StatusBadRequest)
		return
	}

	// Decode the request body into registrationRequest
	var registrationRequest user.UserAuthenticationRequest
	if err := json.NewDecoder(r.Body).Decode(&registrationRequest); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the email format
	if !validator.IsValidEmail(registrationRequest.Email) {
		h.errResp(w, fmt.Errorf("invalid email: %s", registrationRequest.Email), http.StatusBadRequest)
		return
	}

	// Validate the password strength
	if !validator.IsValidPassword(registrationRequest.Password) {
		h.l.Warn("Invalid password provided.")
		h.errResp(w, fmt.Errorf("password must be at least 7 characters long"), http.StatusBadRequest)
		return
	}

	// Generate registration token
	token, err := h.s.Tokens.PrimaryRegistration(registrationRequest.Email, registrationRequest.Password)
	if err != nil {
		switch err {
		case myerrors.ErrDuplicated:
			h.errResp(w, fmt.Errorf("error registering user: already exists"), http.StatusConflict)
		default:
			h.errResp(w, fmt.Errorf("error registering user: invalid request payload: %v", err), http.StatusBadRequest)
		}
		return
	}

	// Send success response
	response := map[string]interface{}{
		"message":     "Confirm your email",
		"token":       token,
		"status_code": http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.l.Error("Failed to send response", zap.Error(err))
	}

	h.l.Debug("User registered successfully", zap.String("email", registrationRequest.Email))
}

// Login authenticates a user and returns an authentication token.
//
// @Summary Login to the system
// @Description Login to the system and get an authentication token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body user.UserAuthenticationRequest true "UserAuthenticationRequest object"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Login attempt initiated...")

	// Check Content-Type header
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		h.errResp(w, fmt.Errorf("invalid Content-Type, expected application/json"), http.StatusBadRequest)
		return
	}

	// Decode the login request payload
	var loginRequest user.UserAuthenticationRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	email := loginRequest.Email
	password := loginRequest.Password

	// Attempt to log in and retrieve the authentication token
	token, err := h.s.Tokens.Login(email, password)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrEmpty), errors.Is(err, myerrors.ErrInvalidCreds):
			h.errResp(w, fmt.Errorf("invalid email or password: %w", err), http.StatusUnauthorized)
		case errors.Is(err, myerrors.ErrInternal), errors.Is(err, myerrors.ErrEmailing):
			h.errResp(w, fmt.Errorf("internal error during login: %w", err), http.StatusInternalServerError)
		}
		return
	}

	// Send success response with token
	response := map[string]interface{}{
		"message":     "Confirm your email",
		"token":       token,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Login successful", zap.String("email", email))
}

func (h *MyHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}
	//! Заставляет задуматься
	userID, ok := utility.GetUserIDFromContext(r.Context())
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

	access, refresh, err := h.s.Tokens.RefreshToken(tmpToken.RefreshToken, userID)
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

// Logout handles user logout and terminates the session.
//
// @Summary Logout the user
// @Description Logs out the user, terminating the session.
// @Tags Auth
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Internal Server Error"
// @Security JWT
// @Router /auth/logout [post]
func (h *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("User logout initiated...")

	// Retrieve the current device ID from the context
	currentDeviceID, ok := utility.GetDeviceIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated: device ID not found"), http.StatusUnauthorized)
		return
	}

	// Retrieve the user ID from the context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated: user ID not found"), http.StatusUnauthorized)
		return
	}

	// Perform logout operation
	if err := h.s.Users.Logout(currentDeviceID, userID); err != nil {
		h.errResp(w, fmt.Errorf("error during logout: %v", err), http.StatusInternalServerError)
		return
	}

	// Send successful logout response
	response := map[string]interface{}{
		"message":     "Logout successful",
		"status_code": http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("User logged out successfully", zap.String("userID", userID))
}

func (h *MyHandler) DeleteTokensHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	deviceID := r.URL.Query().Get("deviceID")
	if (email == "" && deviceID == "") || (email != "" && deviceID != "") {
		jsonresponse.SendErrorResponse(w, errors.New("blank fields and two methods are not allowed"), http.StatusBadRequest)
		return
	}

	err := h.s.Users.DeleteTokens(email, deviceID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":     "Successfuly deleted tokens",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *MyHandler) GetTokenPairsAmountHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("blank fields are not allowed"), http.StatusBadRequest)
		return
	}
	amount, err := h.s.Users.GetTokenPairsAmount(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("error while counting sessions: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":     "Successfuly got ammount",
		"ammount":     amount,
		"status_code": http.StatusOK,
	}
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

	err = h.s.Tokens.ResetPassword(email)
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
	err = h.s.Tokens.ChangePasswordForRecover(email, password, resetToken)
	if err != nil {
		var statusCode = 500
		switch {
		case errors.Is(err, myerrors.ErrEmpty):
			statusCode = http.StatusBadRequest
			break
		case errors.Is(err, myerrors.ErrEmailing):
			statusCode = http.StatusInternalServerError
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
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
