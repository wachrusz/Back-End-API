package v1

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/wachrusz/Back-End-API/pkg/encryption"
	"github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	redisKey = "rate_limit:details"
	totalKey = "rate_limit:total"
)

func (h *MyHandler) ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		contentType = strings.Split(contentType, ";")[0]

		w.Header().Set("Content-Type", contentType)

		next.ServeHTTP(w, r)
	})
}

func (h *MyHandler) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, err := h.getIP(r)
		if err != nil {
			h.errResp(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
			return
		}

		// Создаем контекст с таймаутом для работы с Redis
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Увеличиваем счётчик запросов для данного IP
		newCount, err := h.rdb.HIncrBy(ctx, redisKey, ip, 1).Result()
		if err != nil {
			h.errResp(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
			return
		}

		// Увеличиваем общий счётчик запросов
		_, err = h.rdb.Incr(ctx, totalKey).Result()
		if err != nil {
			h.errResp(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
			return
		}

		// Устанавливаем TTL для хеша, чтобы данные устарели через 1 секунду
		if err := h.rdb.Expire(ctx, redisKey, time.Second).Err(); err != nil {
			h.errResp(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
			return
		}

		// Устанавливаем TTL для общего счётчика
		if err := h.rdb.Expire(ctx, totalKey, time.Second).Err(); err != nil {
			h.errResp(w, fmt.Errorf("internal server error"), http.StatusInternalServerError)
			return
		}

		// Если количество запросов для данного IP больше или равно лимиту, то возвращаем ошибку
		if newCount > h.rateLimitCfg {
			h.errResp(w, fmt.Errorf("Rate-limited. Please, slow down."), http.StatusTooManyRequests)
			return
		}

		// Продолжаем обработку запроса
		next.ServeHTTP(w, r)
	})
}

func (h *MyHandler) getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return ip, nil
}

func (h *MyHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utility.ExtractTokenFromHeader(r)
		if err != nil {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("authentication error: %w", err), http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return encryption.SecretKey, nil
		})

		if err != nil {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("authentification error: %w", err), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid token"), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid token: can't get claims"), http.StatusUnauthorized)
			return
		}
		userIDClaim, ok := claims["sub"]
		if !ok {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid token: no sub in claims"), http.StatusUnauthorized)
			return
		}

		userID, ok := userIDClaim.(string)
		if !ok {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("failed to convert userID into string"), http.StatusUnauthorized)
			return
		}

		expiresClaim, ok := claims["exp"]
		if !ok {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid token: no exp in claims"), http.StatusUnauthorized)
			return
		}

		expiresFloat, ok := expiresClaim.(float64)
		if !ok {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("failed to convert expires into float64"), http.StatusUnauthorized)
			return
		}
		expires := int64(expiresFloat)

		if time.Unix(expires, 0).Before(time.Now()) {
			jsonresponse.SendErrorResponse(w, fmt.Errorf("token expired"), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(setUserIDInContext(r.Context(), userID))
		//h.s.Users.UpdateLastActivity(userID)

		next.ServeHTTP(w, r)
	}
}

func setUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "userID", userID)
}
