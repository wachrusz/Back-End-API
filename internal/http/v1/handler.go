package v1

import (
	"github.com/wachrusz/Back-End-API/internal/service"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"go.uber.org/zap"
	"net/http"
)

type MyHandler struct {
	s *service.Services
	l *zap.Logger
}

func NewHandler(services *service.Services, logger *zap.Logger) *MyHandler {
	return &MyHandler{
		s: services,
		l: logger,
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
