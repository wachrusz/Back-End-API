package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/token"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
	"time"
)

// SendConfirmationEmailTestHandler sends a confirmation email with a code.
//
// @Summary Send confirmation email
// @Description Sends a confirmation email with a generated code.
// @Tags Email
// @Param email query string true "Email address"
// @Param token query string true "Token"
// @Success 200 {string} string "Successfully sent confirmation code"
// @Failure 500 {string} string "Internal server error"
// @Router /email/send-confirmation [post]
func (h *MyHandler) SendConfirmationEmailTestHandler(email, token string, w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Sending confirmation email...")
	confirmationCode, err := utility.GenerateConfirmationCode()
	if err != nil {
		return
	}

	err = h.s.Emails.SaveConfirmationCode(email, confirmationCode, token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully sent confirmation code.",
		"code":        confirmationCode,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetConfirmationCodeTestHandler retrieves the confirmation code for a specific email.
//
// @Summary Get confirmation code
// @Description Retrieves the confirmation code for the provided email.
// @Tags Email
// @Param email query string true "Email address"
// @Success 200 {string} string "Successfully retrieved confirmation code"
// @Failure 400 {string} string "Invalid email"
// @Failure 500 {string} string "Internal server error"
// @Router /email/get-confirmation-code [get]
func (h *MyHandler) GetConfirmationCodeTestHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Retrieving confirmation code...")

	type jsonEmail struct {
		Email string `json:"email"`
	}

	//! ELDER VER
	/*
		var email_struct jsonEmail

			errResp := json.NewDecoder(r.Body).Decode(&email_struct)
			if errResp != nil {
				jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+errResp.Error()), http.StatusBadRequest)
				return
			}
			email := email_struct.Email
	*/

	email := r.URL.Query().Get("email")
	if email == "" {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("Incorrect email"), http.StatusBadRequest)
		return
	}

	code, err := h.s.Emails.GetConfirmationCode(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("Email not found."), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully sent confirmation code.",
		"code":        code,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

// ConfirmEmailHandler confirms the user's email using a confirmation token and code.
//
// @Summary Confirm email
// @Description Confirms the user's email using a token and confirmation code.
// @Tags Email
// @Accept json
// @Produce json
// @Param confirmRequest body token.ConfirmEmailRequest true "Confirmation request"
// @Success 200 {string} string "Successfully confirmed email"
// @Failure 400 {string} string "Invalid request or missing token"
// @Failure 500 {string} string "Internal server error"
// @Router /email/confirm [post]
func (h *MyHandler) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Confirming email...")

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid Content-Type, expected application/json: %v", err), http.StatusBadRequest)
		return
	}

	var confirmRequest token.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("Empty token")
		jsonresponse.SendErrorResponse(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	deviceID, err := utility.GetDeviceIDFromRequest(r)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	token_details, err := h.s.Tokens.ConfirmEmail(token, confirmRequest.EnteredCode, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrInternal) || errors.Is(err, myerrors.ErrEmailing):
			jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
			jsonresponse.SendErrorResponse(w, err, http.StatusUnauthorized)
			break
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":                 "Successfuly confirmed email",
		"token_details":           token_details,
		"access_token_life_time":  time.Minute * 15,
		"refresh_token_life_time": 30 * 24 * time.Hour,
		"status_code":             http.StatusOK,
		"device_id":               deviceID,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

// ConfirmEmailLoginHandler confirms a user's email for login.
//
// @Summary Confirm email for login
// @Description Confirms the user's email for login using a token and confirmation code.
// @Tags Email
// @Accept json
// @Produce json
// @Param confirmRequest body token.ConfirmEmailRequest true "Confirmation request"
// @Success 200 {string} string "Successfully confirmed email for login"
// @Failure 400 {string} string "Invalid request or missing token"
// @Failure 500 {string} string "Internal server error"
// @Router /email/confirm-login [post]
func (h *MyHandler) ConfirmEmailLoginHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Confirming email for login...")
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest token.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("Empty token")
		jsonresponse.SendErrorResponse(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	deviceID, err := utility.GetDeviceIDFromRequest(r)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	token_details, err := h.s.Tokens.ConfirmEmail(token, confirmRequest.EnteredCode, deviceID)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrInternal) || errors.Is(err, myerrors.ErrEmailing):
			jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
			jsonresponse.SendErrorResponse(w, err, http.StatusUnauthorized)
			break
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":                 "Successfuly confirmed email",
		"token_details":           token_details,
		"access_token_life_time":  time.Minute * 15,
		"refresh_token_life_time": 30 * 24 * time.Hour,
		"status_code":             http.StatusOK,
		"device_id":               deviceID,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

// ResetPasswordConfirmHandler confirms the password reset process.
//
// @Summary Confirm password reset
// @Description Confirms the password reset process using a token and code.
// @Tags Password
// @Accept json
// @Produce json
// @Param confirmRequest body token.ConfirmEmailRequest true "Confirmation request"
// @Success 200 {string} string "Successfully confirmed password reset"
// @Failure 400 {string} string "Invalid request or missing token"
// @Failure 500 {string} string "Internal server error"
// @Router /password/reset-confirm [post]
func (h *MyHandler) ResetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Confirming password reset...")
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest token.ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("Empty token")
		jsonresponse.SendErrorResponse(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	err := h.s.Emails.ResetPasswordConfirm(token, confirmRequest.EnteredCode)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrInternal) || errors.Is(err, myerrors.ErrEmailing):
			jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
			jsonresponse.SendErrorResponse(w, err, http.StatusUnauthorized)
			break
		}
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully confirmed email",
		"status_code": http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}
