package v1

import (
	"github.com/wachrusz/Back-End-API/internal/service"
)

type MyHandler struct {
	s *service.Services
}

func NewHandler(s *service.Services) *MyHandler {
	return &MyHandler{s: s}
}
