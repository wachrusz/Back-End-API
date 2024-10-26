package v1

import (
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
)

func (h *MyHandler) SendConfirmationEmailTestHandler(email, token string, w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Sending confirmation email...")
	confirmationCode, err := utility.GenerateConfirmationCode()
	if err != nil {
		return
	}

	err = h.s.Emails.SaveConfirmationCode(email, confirmationCode, token)
	if err != nil {
		h.errResp(w, err, http.StatusInternalServerError)
		return
	}

	response := jsonresponse.CodeResponse{
		Message:    "Confirmation code sent successfully",
		Code:       confirmationCode,
		StatusCode: http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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
				h.errResp(w, errors.New("Invalid request payload: "+errResp.Error()), http.StatusBadRequest)
				return
			}
			email := email_struct.Email
	*/

	email := r.URL.Query().Get("email")
	if email == "" {
		h.errResp(w, fmt.Errorf("Incorrect email"), http.StatusBadRequest)
		return
	}

	code, err := h.s.Emails.GetConfirmationCode(email)
	if err != nil {
		h.errResp(w, fmt.Errorf("Email not found."), http.StatusInternalServerError)
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
