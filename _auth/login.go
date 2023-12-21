//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	enc "backEndAPI/_encryption"
	jsonresponse "backEndAPI/_json_response"
	logger "backEndAPI/_logger"
	"errors"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// @Summary Login to the system
// @Description Login to the system and get an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body auth.UserAuthenticationRequest true "UserAuthenticationRequest object"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}

	var loginRequest UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := loginRequest.Email
	password := loginRequest.Password

	if !checkLoginConds(email, password, w, r) {
		return
	}

	//! SESSIONS

	userID, err_id := getUserIDFromUsersDatabase(email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting session: %v", err_id), http.StatusInternalServerError)
		logger.ErrorLogger.Printf("Unknown exeption in userID %s\n", userID)
	}

	token, err := generateToken(userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Internal Server Error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	deviceID := GetDeviceIDFromRequest(r)

	if IsUserActive(userID) {
		currentUser := ActiveUsers[userID]
		if currentUser.DeviceID == deviceID {
			http.Error(w, fmt.Sprintf("Already logged in"), http.StatusUnauthorized)
			return
		} else {
			removeSessionFromDatabase(currentUser.DeviceID, currentUser.UserID)
			currentUser.DeviceID = deviceID
			ActiveUsers[userID] = currentUser
		}
	}

	//! SAVE SESSIONS
	AddActiveUser(userID, email, deviceID, token)

	saveSessionToDatabase(email, deviceID, userID, token)

	response := map[string]interface{}{
		"message":       "Successfuly logged in",
		"status_code":   http.StatusOK,
		"token":         token,
		"refresh_token": "temp_blank",
	}
	json.NewEncoder(w).Encode(response)
}

func checkLoginConds(email, password string, w http.ResponseWriter, r *http.Request) bool {

	if email == "" || password == "" {
		jsonresponse.SendErrorResponse(w, errors.New("Missing email or password: "), http.StatusBadRequest)
		logger.ErrorLogger.Printf("Missing email or password in login request from %s\n", r.RemoteAddr)
		return false
	}

	if !isValidCredentials(email, password) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email or password: "), http.StatusUnauthorized)
		logger.ErrorLogger.Printf("Invalid email or password in login request from %s\n", r.RemoteAddr)
		return false
	}
	return true
}

func generateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	})

	tokenString, err := token.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
