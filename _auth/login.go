//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	enc "backEndAPI/_encryption"
	jsonresponse "backEndAPI/_json_response"
	logger "backEndAPI/_logger"
	mydb "backEndAPI/_mydatabase"
	"errors"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

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

	token, err := generateToken(userID, time.Minute*15)
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
	AddActiveUser(userID, email, deviceID, token.AccessToken)

	saveSessionToDatabase(email, deviceID, userID, token.AccessToken)

	response := map[string]interface{}{
		"message":       "Successfuly logged in",
		"token_details": token,
		"status_code":   http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

// *NEW
// ! СРОЧНО ДОДЕЛАТЬ ЭТО
func RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}
	//! Заставляет задуматься
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	//! Сомнительно
	type token struct {
		RefreshToken string `json:"refresh_token"`
	}

	var tmp_token token

	err := json.NewDecoder(r.Body).Decode(&tmp_token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	tokenDetails, err := refreshToken(tmp_token.RefreshToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	err = updateTokenInDB(userID, tokenDetails.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update token in DB: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":       "Successfully logged in",
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
		"status_code":   http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// *NEW
// ! Это выглядит подозрительно плохо, есть вероятность что нахуй буду послан при попытке запустить(так и было)
func updateTokenInDB(userID, newAccessToken string) error {
	encryptedToken, err := enc.EncryptToken(newAccessToken)
	if err != nil {
		return err
	}

	SetAccessToken(userID, newAccessToken)

	query := `
		UPDATE sessions
		SET token = $1,
		expires_at = NOW() + INTERVAL '15 minutes'
		WHERE user_id = $2;
	`
	_, err = mydb.GlobalDB.Exec(query, encryptedToken, userID)
	return err
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

func generateToken(userID string, duration time.Duration) (*TokenDetails, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenExpiresAt := time.Now().Add(30 * 24 * time.Hour).Unix()

	return &TokenDetails{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    refreshTokenExpiresAt,
	}, nil
}

func refreshToken(refreshTokenString string) (*TokenDetails, error) {
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := refreshToken.Claims.(jwt.Claims); !ok && !refreshToken.Valid {
		return nil, fmt.Errorf("Invalid refresh token")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Failed to parse refresh token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("Failed to get user ID from refresh token")
	}

	return generateToken(userID, time.Minute*15)
}
