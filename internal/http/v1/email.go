package v1

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/token"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"

	"net/http"

	"github.com/wachrusz/Back-End-API/pkg/validator"
	"go.uber.org/zap"
)

// ConfirmResponse решено вынести из пакета jsonresponse во избежание циклических зависимостей, так как требует token.Details.
type ConfirmResponse struct {
	Message              string        `json:"message"`
	TokenDetails         token.Details `json:"token_details"`
	AccessTokenLifeTime  int64         `json:"access_token_life_time"`
	RefreshTokenLifeTime int64         `json:"refresh_token_life_time"`
	StatusCode           int           `json:"status_code"`
	DeviceId             string        `json:"device_id"`
}

// ConfirmEmailRegisterHandler confirms the user's email using a confirmation RefreshToken and code during registration.
//
// @Summary Confirm email
// @Description Confirms the user's email using a RefreshToken and confirmation code during registration.
// @Tags Auth
// @Accept json
// @Produce json
// @Param confirmRequest body token.ConfirmEmailRequest true "Confirmation request"
// @Success 200 {object} ConfirmResponse "Successfully confirmed email"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request or missing RefreshToken"
// @Failure 401 {object} jsonresponse.CodeError "Invalid code"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal server error"
// @Router /auth/register/confirm [post]
func (h *MyHandler) ConfirmEmailRegisterHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Confirming email...")

	var confirmRequest token.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		h.errResp(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("empty RefreshToken")
		h.errResp(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	deviceID, err := utility.GetDeviceIDFromRequest(r)
	if err != nil {
		h.errResp(w, fmt.Errorf("internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	details, err := h.s.Tokens.ConfirmEmailRegister(token, confirmRequest.EnteredCode, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrInternal) || errors.Is(err, myerrors.ErrEmailing):
			h.errResp(w, err, http.StatusInternalServerError)
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
			h.errResp(w, err, http.StatusBadRequest)
			break
		case errors.Is(err, myerrors.ErrCode) || errors.Is(err, myerrors.ErrLocked):
			h.errAuthResp(w, err, details.RemainingAttempts, details.LockDuration, http.StatusUnauthorized)
			break
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := ConfirmResponse{
		Message:              "Successfully confirmed email",
		TokenDetails:         details,
		AccessTokenLifeTime:  60 * 15,
		RefreshTokenLifeTime: 30 * 24 * 60 * 60,
		StatusCode:           http.StatusOK,
		DeviceId:             deviceID,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// ConfirmEmailLoginHandler confirms a user's email for login.
//
// @Summary Confirm email for login
// @Description Confirms the user's email for login using a RefreshToken and confirmation code.
// @Tags Auth
// @Accept json
// @Produce json
// @Param confirmRequest body token.ConfirmEmailRequest true "Confirmation request"
// @Success 200 {object} ConfirmResponse             "Successfully confirmed email for login"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request or missing RefreshToken"
// @Failure 401 {object} jsonresponse.CodeError     "Invalid code"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal server error"
// @Router /auth/login/confirm [post]
func (h *MyHandler) ConfirmEmailLoginHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Confirming email for login...")

	var confirmRequest token.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		h.errResp(w, fmt.Errorf("empty refresh token"), http.StatusBadRequest)
		return
	}

	deviceID, err := utility.GetDeviceIDFromRequest(r)
	if err != nil {
		h.errResp(w, fmt.Errorf("internal Server Error: %v", err), http.StatusBadRequest)
		return
	}

	details, err := h.s.Tokens.ConfirmEmailLogin(token, confirmRequest.EnteredCode, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrInternal) || errors.Is(err, myerrors.ErrEmailing):
			h.errResp(w, err, http.StatusInternalServerError)
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
			h.errResp(w, err, http.StatusBadRequest)
			break
		case errors.Is(err, myerrors.ErrCode) || errors.Is(err, myerrors.ErrLocked):
			h.errAuthResp(w, err, details.RemainingAttempts, details.LockDuration, http.StatusUnauthorized)
			break
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := ConfirmResponse{
		Message:              "Successfully confirmed email",
		TokenDetails:         details,
		AccessTokenLifeTime:  60 * 15,
		RefreshTokenLifeTime: 30 * 24 * 60 * 60,
		StatusCode:           http.StatusOK,
		DeviceId:             deviceID,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// ConfirmResetResponse решено вынести из пакета jsonresponse во избежание циклических зависимостей, так как требует token.Details.
type ConfirmResetResponse struct {
	Message            string                  `json:"message"`
	TokenDetails       token.ResetTokenDetails `json:"token_details"`
	ResetTokenLifeTime int64                   `json:"reset_token_life_time"`
	StatusCode         int                     `json:"status_code"`
	DeviceId           string                  `json:"device_id"`
}

// ResetPasswordConfirmHandler confirms the password reset process.
//
// @Summary Confirm password reset
// @Description Confirms the password reset process using a RefreshToken and code.
// @Tags Auth
// @Accept json
// @Produce json
// @Param confirmRequest body token.ConfirmEmailRequest true 	"Confirmation request"
// @Success 200 {object} jsonresponse.SuccessResponse 			"Successfully confirmed password reset"
// @Failure 400 {object} jsonresponse.ErrorResponse 			"Unauthorized"
// @Failure 401 {object} jsonresponse.CodeError 				"Unauthorized with remaining attempts"
// @Failure 500 {object} jsonresponse.ErrorResponse 			"Internal server error"
// @Router /auth/password/confirm [post]
func (h *MyHandler) ResetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Confirming password reset...")
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		h.errResp(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest token.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		h.errResp(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("empty RefreshToken")
		h.errResp(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	deviceID, err := utility.GetDeviceIDFromRequest(r)
	if err != nil {
		h.errResp(w, fmt.Errorf("internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	details, err := h.s.Tokens.ResetPasswordConfirm(token, confirmRequest.EnteredCode, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrInternal) || errors.Is(err, myerrors.ErrEmailing):
			h.errResp(w, err, http.StatusInternalServerError)
			break
		case errors.Is(err, myerrors.ErrInvalidToken) || errors.Is(err, myerrors.ErrExpiredCode):
			h.errResp(w, err, http.StatusBadRequest)
			break
		case errors.Is(err, myerrors.ErrCode) || errors.Is(err, myerrors.ErrLocked):
			h.errAuthResp(w, err, details.RemainingAttempts, details.LockDuration, http.StatusUnauthorized)
			break
		}
		return
	}

	response := ConfirmResetResponse{
		Message:            "Successfully confirmed email",
		TokenDetails:       details,
		ResetTokenLifeTime: 60 * 15,
		StatusCode:         http.StatusOK,
		DeviceId:           deviceID,
	}

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// RegisterUserHandler register a new user and send confirmation email.
//
// @Summary Register user
// @Description Register a new user and send confirmation email.
// @Tags Auth
// @Accept json
// @Produce json
// @Param confirmRequest body UserAuthenticationRequest true 	"Credentials"
// @Success 200 {object} jsonresponse.TokenResponse "User registered successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 409 {object} jsonresponse.ErrorResponse "User already exists"
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
	var registrationRequest UserAuthenticationRequest
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

	// Generate registration RefreshToken
	token, err := h.s.Tokens.Register(registrationRequest.Email, registrationRequest.Password)
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
	response := jsonresponse.TokenResponse{
		Message:    "Confirm your email",
		Token:      token,
		StatusCode: http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.l.Error("Failed to send response", zap.Error(err))
	}

	h.l.Debug("User registered successfully", zap.String("email", registrationRequest.Email))
}

// LoginUserHandler authenticates a user and returns an authentication RefreshToken.
//
// @Summary LoginUserHandler to the system
// @Description LoginUserHandler to the system and get an authentication RefreshToken.
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body UserAuthenticationRequest true "UserAuthenticationRequest object"
// @Success 200 {object} jsonresponse.TokenResponse "LoginUserHandler successful"
// @Failure 400 {object} jsonresponse.ErrorResponse "Bad Request"
// @Failure 401 {object} jsonresponse.ErrorResponse "Unauthorized"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (h *MyHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Login attempt initiated...")

	// Check Content-Type header
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		h.errResp(w, fmt.Errorf("invalid Content-Type, expected application/json"), http.StatusBadRequest)
		return
	}

	// Decode the login request payload
	var loginRequest UserAuthenticationRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	email := loginRequest.Email
	password := loginRequest.Password

	// Attempt to log in and retrieve the authentication RefreshToken
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

	// Send success response with RefreshToken
	response := jsonresponse.TokenResponse{
		Message:    "Confirm your email",
		Token:      token,
		StatusCode: http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Login successful", zap.String("email", email))
}

// ResetPasswordHandler sends a password reset token to the user's email.
//
// @Summary Reset password
// @Description This endpoint allows users to request a password reset by providing their email. If the email is valid, a reset token will be sent to it.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   body  body  user.ResetPasswordRequest  true  "Reset password request with email"
// @Success 200 {object} jsonresponse.SuccessResponse "Successfully sent email with reset password token"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid Content-Type or invalid email"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error"
// @Router /auth/password [post]
func (h *MyHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Attempting to reset password...")
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("empty 'Content-Type' HEADER")
		h.errResp(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var resetRequest user.ResetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		h.errResp(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	if !validator.IsValidEmail(email) {
		h.errResp(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	token, err := h.s.Tokens.ResetPassword(email)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrInvalidCreds):
			h.errResp(w, err, http.StatusBadRequest)
			break
		case errors.Is(err, myerrors.ErrEmailing) || errors.Is(err, myerrors.ErrInternal):
			h.errResp(w, err, http.StatusInternalServerError)
			break
		}
		return
	}

	response := jsonresponse.TokenResponse{
		Message:    "Successfully sent email with reset password token",
		Token:      token,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// UserAuthenticationRequest is for auth requests
type UserAuthenticationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
