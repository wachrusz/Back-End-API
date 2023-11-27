//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	logger "backEndAPI/_logger"
	user "backEndAPI/_user"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	secretAPIKey          string
	sessionMutex          sync.Mutex
	globalCurrentDeviceID string
)

const (
	SessionName = "session-name"
	deviceIDKey = "deviceID"
)

func setUserIDInContext(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, "Key", userID)
}

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/auth/login", Login).Methods("POST")
	router.HandleFunc("/auth/logout", Logout).Methods("POST")
}

func setGlobalCurrentDeviceID(deviceID string) {
	globalCurrentDeviceID = deviceID
	log.Print("setGlobalCurrentDeviceID: ", globalCurrentDeviceID)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := GetDeviceIDFromRequest(r)
		userID, err := GetUserIDFromSessionDatabase(deviceID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting session: %v", err), http.StatusInternalServerError)
			logger.ErrorLogger.Printf("Unknown exeption in userID %s\n", userID)
			return
		}
		if !IsUserActive(userID) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		intUserID, _ := strconv.Atoi(userID)
		r = r.WithContext(setUserIDInContext(r.Context(), intUserID))

		next.ServeHTTP(w, r)
	}
}

func SetAPIKey() (string, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return "", err
	}

	log.Println(apiKey)

	return apiKey, nil
}

func generateAPIKey() (string, error) {
	// Генерируем случайные байты для ключа
	keyBytes := make([]byte, 32)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", err
	}

	// Преобразуем байты в строку шестнадцатеричного формата
	apiKey := hex.EncodeToString(keyBytes)

	return apiKey, nil
}

func GetAPIKey() string {
	return secretAPIKey
}

func isValidCredentials(username, password string) bool {
	hashedPassword, ok := user.GetHashedPasswordByUsername(username)
	if ok != nil {
		return false
	}
	if comparePasswords(hashedPassword, password) {
		return true
	}
	return false
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
