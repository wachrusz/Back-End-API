//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	logger "backEndAPI/_logger"
	"encoding/json"

	"github.com/dgrijalva/jwt-go"

	"fmt"
	"net/http"
)

// @Summary Login to the system
// @Description Login to the system and get an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {

	var (
		username string
		password string
		userID   string
	)

	username = r.FormValue("username")
	password = r.FormValue("password")
	apiKey := r.Header.Get("API-Key")

	if !checkLoginConds(apiKey, username, password, w, r) {
		return
	}

	token, err := generateToken(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//! SESSIONS

	userID, err_id := getUserIDFromUsersDatabase(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting session: %v", err_id), http.StatusInternalServerError)
		logger.ErrorLogger.Printf("Unknown exeption in userID %s\n", userID)
	}

	deviceID := GetDeviceIDFromRequest(r)

	if IsUserActive(userID) { //in session
		currentUser := activeUsers[userID]
		if currentUser.DeviceID == deviceID {
			http.Error(w, fmt.Sprintf("Already logged in"), http.StatusUnauthorized)
			return
		} else {
			currentUser.DeviceID = deviceID
			activeUsers[userID] = currentUser
			removeSessionFromDatabase(deviceID)
		}
	}

	//! SAVE SESSIONS
	AddActiveUser(userID, username, deviceID)

	saveSessionToDatabase(username, deviceID, userID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
	json.NewEncoder(w).Encode(map[string]string{"token": token})

	logger.InfoLogger.Printf("User %s logged in from %s\n", username, r.RemoteAddr)
}

func checkLoginConds(apiKey, username, password string, w http.ResponseWriter, r *http.Request) bool {
	if apiKey != secretAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API key"))
		return false
	}

	if username == "" || password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		logger.ErrorLogger.Printf("Missing username or password in login request from %s\n", r.RemoteAddr)
		return false
	}

	if !isValidCredentials(username, password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		logger.ErrorLogger.Printf("Invalid username or password in login request from %s\n", r.RemoteAddr)
		return false
	}
	return true
}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	tokenString, err := token.SignedString([]byte(GetAPIKey()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
