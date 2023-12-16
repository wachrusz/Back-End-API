//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	enc "backEndAPI/_encryption"
	logger "backEndAPI/_logger"

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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Content-Type, expected application/json"))
		return
	}

	var loginRequest UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request payload"))
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
	json.NewEncoder(w).Encode(map[string]string{"token": token})

	logger.InfoLogger.Printf("User %s logged in from %s\n", email, r.RemoteAddr)
}

func checkLoginConds(email, password string, w http.ResponseWriter, r *http.Request) bool {

	if email == "" || password == "" {
		http.Error(w, "Missing email or password", http.StatusBadRequest)
		logger.ErrorLogger.Printf("Missing email or password in login request from %s\n", r.RemoteAddr)
		return false
	}

	if !isValidCredentials(email, password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
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
