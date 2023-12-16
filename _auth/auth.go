//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	"context"
	"net/http"
	"sync"

	enc "backEndAPI/_encryption"
	logger "backEndAPI/_logger"
	user "backEndAPI/_user"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Struct for auth requests
type UserAuthenticationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	//secretAPIKey          string = "61501ebfe5eec2610a486da6da176bb810bbac93d7e5fb928545a4a695c7532d" //! ;< bro died вiд кринжу
	sessionMutex sync.Mutex
)

const (
	SessionName = "session-name"
	deviceIDKey = "deviceID"
)

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/auth/login", Login).Methods("POST")
	router.HandleFunc("/auth/logout", Logout).Methods("POST")
	router.HandleFunc("/auth/register", RegisterUser).Methods("POST")
	router.HandleFunc("/auth/login/vk", handleVKLogin).Methods("POST")
	router.HandleFunc("/auth/login/google", handleGoogleLogin).Methods("POST")
	router.HandleFunc("/auth/register/confirm-email", ConfirmEmailHandler).Methods("POST")
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			logger.ErrorLogger.Printf("Error in tokenString")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(enc.SecretKey), nil
		})

		if err != nil {
			logger.ErrorLogger.Printf("Error in parsing")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			logger.ErrorLogger.Printf("Invalid token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.ErrorLogger.Printf("Claims error")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userIDClaim, ok := claims["sub"]
		if !ok {
			logger.ErrorLogger.Printf("No 'sub' claim in token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := userIDClaim.(string)
		if !ok {
			logger.ErrorLogger.Printf("Failed to convert 'sub' claim to string")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !IsUserActive(userID) {
			logger.ErrorLogger.Printf("Not active, userID: %v", userID)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(setUserIDInContext(r.Context(), userID))

		next.ServeHTTP(w, r)
	}
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

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}
