package auth

import (
	email_conf "backEndAPI/_email"
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Content-Type, expected application/json"))
		return
	}

	var registrationRequest UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&registrationRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
		return
	}

	email := registrationRequest.Email
	password := registrationRequest.Password
	if !isValidEmail(email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid email"))
		return
	}

	token, err := utility.GenerateRegisterJWTToken(email, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating confirmation token"))
		return
	}

	err = email_conf.SendConfirmationEmail(email, token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	response := map[string]interface{}{
		"message": "Confirm your email",
		"token":   token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
