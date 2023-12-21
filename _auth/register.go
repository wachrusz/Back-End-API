package auth

import (
	email_conf "backEndAPI/_email"
	jsonresponse "backEndAPI/_json_response"
	"errors"

	//user "backEndAPI/_user"
	utility "backEndAPI/_utility"

	"encoding/json"
	"net/http"
	"net/mail"
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
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var registrationRequest UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&registrationRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := registrationRequest.Email
	password := registrationRequest.Password
	if !isValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	token, err := utility.GenerateRegisterJWTToken(email, password)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error generating confirmation token: "+err.Error()), http.StatusInternalServerError)
		return
	}

	err = email_conf.SendConfirmationEmail(email, token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error sending confirm email: "+err.Error()), http.StatusInternalServerError)
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

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var resetRequest ResetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	if !isValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	token, err := utility.GenerateResetJWTToken(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error generating confirmation token: "+err.Error()), http.StatusInternalServerError)
		return
	}

	err = email_conf.SendConfirmationEmail(email, token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error sending confirm email: "+err.Error()), http.StatusInternalServerError)
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

func ChangePasswordForRecoverHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}
	var registrationRequest UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&registrationRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := registrationRequest.Email
	password := registrationRequest.Password
	if !isValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	resetToken := r.Header.Get("Authorization")
	if resetToken == "" {
		jsonresponse.SendErrorResponse(w, errors.New("Reset token is required"), http.StatusBadRequest)
		return
	}
	_, err = utility.VerifyResetJWTToken(resetToken)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired reset token: "+err.Error()), http.StatusUnauthorized)
		return
	}

	err = ResetPassword(email, password)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error resetting password: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfuly reseted password",
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	Login(w, r)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
