package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/httprate"
	"github.com/wachrusz/Back-End-API/pkg/encryption"
	"github.com/wachrusz/Back-End-API/pkg/json_response"
	"net/http"
	"strings"
	"time"
)

func (h *MyHandler) ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		contentType = strings.Split(contentType, ";")[0]

		w.Header().Set("Content-Type", contentType)

		next.ServeHTTP(w, r)
	})
}

func (h *MyHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			err := errors.New("Error in tokenString")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(encryption.SecretKey), nil
		})
		if err != nil {
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			err := errors.New("Invalid RefreshToken")
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
			err := errors.New("No 'sub' claim in RefreshToken")
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

		if !h.s.Users.IsUserActive(userID) {
			err := errors.New("Inactive user")
			jsonresponse.SendErrorResponse(w, errors.New("Unauthorized: "+err.Error()), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(setUserIDInContext(r.Context(), userID))
		r = r.WithContext(setDeviceIDInContext(r.Context(), deviceID))
		h.s.Users.UpdateLastActivity(userID)

		next.ServeHTTP(w, r)
	}
}

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}
func setDeviceIDInContext(ctx context.Context, deviceID string) context.Context {
	return context.WithValue(ctx, "device_id", deviceID)
}

func (h *MyHandler) RateLimitMiddleware(next http.Handler) http.Handler {
	rateLimiter := httprate.Limit(
		h.rateLimit,
		time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			h.errResp(w, fmt.Errorf("Rate-limited. Please, slow down."), http.StatusTooManyRequests)
		}),
	)

	return rateLimiter(next)
}
