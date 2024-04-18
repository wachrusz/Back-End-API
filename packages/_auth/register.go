package auth

import (
	"database/sql"
	"errors"
	"fmt"

	service "main/packages/_auth/service"
	email_conf "main/packages/_email"
	jsonresponse "main/packages/_json_response"
	mydb "main/packages/_mydatabase"

	//user "../_user"
	utility "main/packages/_utility"

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

	err, used := isEmailUsed(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if used {
		jsonresponse.SendErrorResponse(w, errors.New("Email already exists: "), http.StatusBadRequest)
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

	type UserPasswordReset struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		ResetToken string `json:"reset_token"`
	}

	var resetRequest UserPasswordReset

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	password := resetRequest.Password
	resetToken := resetRequest.ResetToken
	if !isValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	if resetToken == "" {
		jsonresponse.SendErrorResponse(w, errors.New("Reset token is required"), http.StatusBadRequest)
		return
	}
	_, err = utility.VerifyResetJWTToken(resetToken)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired reset token: "+err.Error()), http.StatusUnauthorized)
		return
	}
	claims, err := utility.ParseResetToken(resetToken)
	if claims["code_used"].(bool) {
		jsonresponse.SendErrorResponse(w, errors.New("Token has already been used: "+err.Error()), http.StatusUnauthorized)
		return
	} else {
		claims["code_used"] = true
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

	userID, err := service.GetUserIDFromUsersDatabase(email)
	err = InvalidateTokensByUserID(userID)
	if err != nil && err != sql.ErrNoRows {
		response := map[string]interface{}{
			"message":     "Error appeared in token_ivalidation",
			"status_code": http.StatusInternalServerError,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusFound)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// *NEW
func isEmailUsed(email string) (error, bool) {

	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	err := mydb.GlobalDB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return fmt.Errorf("Error getting email: %v", err), false
	}

	return nil, count > 0
}
