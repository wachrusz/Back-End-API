package v1

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/wachrusz/Back-End-API/internal/repository"
	"github.com/wachrusz/Back-End-API/internal/service"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"go.uber.org/zap"
	"net/http"
)

type MyHandler struct {
	s            *service.Services
	l            *zap.Logger
	m            *repository.Models
	rdb          *redis.Client
	rateLimitCfg int64
}

func NewHandler(services *service.Services, logger *zap.Logger, models *repository.Models, cache *redis.Client, rateLimit int64) *MyHandler {
	if rateLimit <= 0 {
		rateLimit = 10
	}

	return &MyHandler{
		s:            services,
		l:            logger,
		m:            models,
		rateLimitCfg: rateLimit,
		rdb:          cache,
	}
}

// errResp logs all errors and sends json response to rhe server
func (h *MyHandler) errResp(w http.ResponseWriter, err error, statusCode int) {
	h.l.Error("Error occurred",
		zap.String("error", err.Error()),
		zap.Int("status_code", statusCode),
	)
	jsonresponse.SendErrorResponse(w, err, statusCode)
}

// errAuthResp logs all authentication errors and sends json response to rhe server
func (h *MyHandler) errAuthResp(w http.ResponseWriter, err error, attempts, lockedDuration, statusCode int) {
	h.l.Error("Error occurred",
		zap.String("error", err.Error()),
		zap.Int("attempts", attempts),
		zap.Int("locked_duration", lockedDuration),
		zap.Int("status_code", statusCode),
	)

	r := jsonresponse.CodeError{
		Error:        err.Error(),
		Attempts:     attempts,
		LockDuration: lockedDuration,
		StatusCode:   statusCode,
	}

	jsonData, err := json.Marshal(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}
