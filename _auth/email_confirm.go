//go:build !exclude_swagger
// +build !exclude_swagger

// Package profile provides profile information and it's functionality.
package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	jsonresponse "backEndAPI/_json_response"
	mydb "backEndAPI/_mydatabase"
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
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest ConfirmEmailRequest
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

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err := utility.VerifyRegisterJWTToken(token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired token: "+err.Error()), http.StatusUnauthorized)
		return
	}

	if !email_conf.CheckConfirmationCode(registerRequest.Email, confirmRequest.EnteredCode) {
		err := errors.New("Error in CheckConfirmationCode")
		jsonresponse.SendErrorResponse(w, errors.New("Wrong confiramtion code: "+err.Error()), http.StatusInternalServerError)
		return
	}

	err = email_conf.ConfirmEmail(registerRequest.Email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error confirming email: "+err.Error()), http.StatusInternalServerError)
		return
	}

	err = user.RegisterUser(registerRequest.Email, registerRequest.Password)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error registring user: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully confirmed email",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
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
// @Router /auth/login/reset-confirm [post]
func ResetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest ConfirmEmailRequest
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

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err := utility.VerifyResetJWTToken(token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired token: "+err.Error()), http.StatusUnauthorized)
		return
	}

	if !email_conf.CheckConfirmationCode(registerRequest.Email, confirmRequest.EnteredCode) {
		err := errors.New("Error in CheckConfirmationCode")
		jsonresponse.SendErrorResponse(w, errors.New("Wrong confiramtion code: "+err.Error()), http.StatusInternalServerError)
		return
	}

	err = email_conf.ConfirmEmail(registerRequest.Email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error confirming email: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully confirmed email",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func ResetPassword(email, password string) error {
	hashedPassword, err := user.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = mydb.GlobalDB.Exec("UPDATE users SET hashed_password = $1 WHERE email = $2", hashedPassword, email)
	return err
}
