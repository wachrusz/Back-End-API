//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	service2 "github.com/wachrusz/Back-End-API/internal/auth/service"
	email_conf "github.com/wachrusz/Back-End-API/internal/email"
	"github.com/wachrusz/Back-End-API/internal/user"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Struct for auth requests
type UserAuthenticationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
// secretAPIKey          string = "" //! ;< bro died вiд кринжу
// sessionMutex sync.Mutex //! його теж дуже шкода
)

const (
	SessionName = "session-name"
	deviceIDKey = "deviceID"
)

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}
func setDeviceIDInContext(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, "device_id", deviceID)
}

/*
func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/auth/login", Login).Methods("POST")
	//* NEW
	router.HandleFunc("/auth/login/confirm", ConfirmEmailLoginHandler).Methods("POST")

	router.HandleFunc("/auth/logout", AuthMiddleware(Logout)).Methods("POST")
	router.HandleFunc("/auth/register", RegisterUser).Methods("POST")
	router.HandleFunc("/auth/login/vk", handleVKLogin).Methods("GET")
	router.HandleFunc("/auth/login/google", handleGoogleLogin).Methods("GET")
	router.HandleFunc("/auth/register/confirm-email", ConfirmEmailHandler).Methods("POST")
	router.HandleFunc("/auth/login/reset/password/confirm", ResetPasswordConfirmHandler).Methods("POST")
	router.HandleFunc("/auth/login/reset/password", ResetPasswordHandler).Methods("POST")
	router.HandleFunc("/auth/login/reset/password/put", ChangePasswordForRecoverHandler).Methods("PUT")
	//* Fixed
	router.HandleFunc("/auth/refresh", AuthMiddleware(RefreshTokenHandler)).Methods("POST")

	router.HandleFunc("/dev/confirmation-code/get", email_conf.GetConfirmationCodeTestHandler).Methods("GET")

	router.HandleFunc("/auth/tokens/delete", service2.DeleteTokensHandler).Methods("DELETE")
	router.HandleFunc("/auth/tokens/ammount", service2.GetTokenPairsAmmountHandler).Methods("GET")
}
*/

func RegisterHandlers(router chi.Router) {
	// Auth routes
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", Login)
		r.Post("/login/confirm", ConfirmEmailLoginHandler)
		r.Post("/logout", AuthMiddleware(Logout))
		r.Post("/register", RegisterUser)
		r.Post("/register/confirm-email", ConfirmEmailHandler)

		// Password reset routes
		r.Route("/login/reset", func(r chi.Router) {
			r.Post("/password", ResetPasswordHandler)
			r.Post("/password/confirm", ResetPasswordConfirmHandler)
			r.Put("/password/put", ChangePasswordForRecoverHandler)
		})

		// Token routes
		r.Post("/refresh", AuthMiddleware(RefreshTokenHandler))
		r.Delete("/tokens/delete", service2.DeleteTokensHandler)
		r.Get("/tokens/amount", service2.GetTokenPairsAmmountHandler)
	})

	// OAuth login routes
	router.Route("/auth/login", func(r chi.Router) {
		r.Get("/vk", handleVKLogin)
		r.Get("/google", handleGoogleLogin)
	})

	// Developer routes
	router.Get("/dev/confirmation-code/get", email_conf.GetConfirmationCodeTestHandler)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			err := errors.New("Error in tokenString")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(enc.SecretKey), nil
		})
		if err != nil {
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			err := errors.New("Invalid token")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			err := errors.New("Claims error")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}
		userIDClaim, ok := claims["sub"]
		if !ok {
			err := errors.New("No 'sub' claim in token")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		userID, ok := userIDClaim.(string)
		if !ok {
			err := errors.New("Failed to convert 'sub' claim to string")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		deviceID, ok := claims["device_id"].(string)
		if !ok {
			err := errors.New("Failed to convert 'sub' claim to string")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		if !service2.IsUserActive(userID) {
			err := errors.New("Inactive user")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(setUserIDInContext(r.Context(), userID))
		r = r.WithContext(setDeviceIDInContext(r.Context(), deviceID))
		service2.UpdateLastActivity(userID)

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

func GetDeviceIDFromContext(ctx context.Context) (string, bool) {
	deviceID, ok := ctx.Value("device_id").(string)
	return deviceID, ok
}
