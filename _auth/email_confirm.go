//go:build !exclude_swagger
// +build !exclude_swagger

// Package profile provides profile information and it's functionality.
package auth

import (
	"encoding/json"
	"net/http"

	utility "backEndAPI/_utility"

	email_conf "backEndAPI/_email"
	//utility "backEndAPI/_utility"
	user "backEndAPI/_user"
)

// ConfirmEmailRequest структура представляет запрос на подтверждение электронной почты.
type ConfirmEmailRequest struct {
	Token       string `json:"token"`
	EnteredCode string `json:"code"`
}

// @Summary Confirm user email
// @Description Confirm the user's email using the provided token and confirmation code.
// @Tags Auth
// @Produce json
// @Param confirmEmailRequest body ConfirmEmailRequest true "Confirm Email Request"
// @Success 200 {string} string "Email confirmed successfully"
// @Failure 400 {string} string "Invalid request payload or Content-Type"
// @Failure 401 {string} string "Invalid or expired token"
// @Failure 500 {string} string "Error confirming email or registering user"
// @Router /auth/register/confirm-email [post]
func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Content-Type, expected application/json"))
		return
	}

	var confirmRequest ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return
	}

	token := confirmRequest.Token
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Token is required"))
		return
	}

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err := utility.VerifyRegisterJWTToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid or expired token"))
		return
	}

	if !email_conf.CheckConfirmationCode(registerRequest.Email, confirmRequest.EnteredCode) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Wrong confiramtion code"))
		return
	}

	err = email_conf.ConfirmEmail(registerRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error confirming email"))
		return
	}

	err = user.RegisterUser(registerRequest.Email, registerRequest.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error registring user"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email confirmed successfully"))
}

// @Summary Confirm user email for password reset
// @Description Confirm the user's email using the provided token and confirmation code.
// @Tags Auth
// @Produce json
// @Param confirmEmailRequest body ConfirmEmailRequest true "Confirm Email Request"
// @Success 200 {string} string "Email confirmed successfully"
// @Failure 400 {string} string "Invalid request payload or Content-Type"
// @Failure 401 {string} string "Invalid or expired token"
// @Failure 500 {string} string "Error confirming email or reseting password"
// @Router /auth/login/reset [post]
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Content-Type, expected application/json"))
		return
	}

	var confirmRequest ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return
	}

	token := confirmRequest.Token
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Token is required"))
		return
	}

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err := utility.VerifyRegisterJWTToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid or expired token"))
		return
	}

	if !email_conf.CheckConfirmationCode(registerRequest.Email, confirmRequest.EnteredCode) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Wrong confiramtion code"))
		return
	}

	err = email_conf.ConfirmEmail(registerRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error confirming email"))
		return
	}
	/*
		err = ResetPassword(registerRequest.Email, registerRequest.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error registring user"))
			return
		}
	*/
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email confirmed successfully"))
}
