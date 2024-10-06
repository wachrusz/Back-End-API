package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/logger"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
	"time"
)

func (h *MyHandler) SendConfirmationEmailTestHandler(email, token string, w http.ResponseWriter, r *http.Request) {
	confirmationCode, err := utility.GenerateConfirmationCode()
	if err != nil {
		logger.ErrorLogger.Printf("Error in generating confirmation code for Email: %v", email)
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

func (h *MyHandler) GetConfirmationCodeTestHandler(w http.ResponseWriter, r *http.Request) {
	type jsonEmail struct {
		Email string `json:"email"`
	}

	//! ELDER VER
	/*
		var email_struct jsonEmail

			err := json.NewDecoder(r.Body).Decode(&email_struct)
			if err != nil {
				jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
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

func (h *MyHandler) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid Content-Type, expected application/json: %v", err), http.StatusBadRequest)
		return
	}

	var confirmRequest email.ConfirmEmailRequest
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

	deviceID, err := h.getDeviceIDFromRequest(r)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	token_details, err := h.s.Emails.ConfirmEmail(token, confirmRequest.EnteredCode, deviceID)
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

func (h *MyHandler) ConfirmEmailLoginHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest email.ConfirmEmailRequest
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

	deviceID, err := h.getDeviceIDFromRequest(r)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("internal Server Error: %v", err), http.StatusInternalServerError)
		return
	}

	token_details, err := h.s.Emails.ConfirmEmail(token, confirmRequest.EnteredCode, deviceID)
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

func (h *MyHandler) ResetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest email.ConfirmEmailRequest
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
