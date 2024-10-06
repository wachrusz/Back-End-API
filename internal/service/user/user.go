package user

import (
	"github.com/wachrusz/Back-End-API/internal/mydatabase"
	"sync"
)

type Service struct {
	repo        *mydatabase.Database
	ActiveUsers map[string]ActiveUser
	mutex       sync.Mutex
	activeMu    sync.Mutex
}

func NewService(repo *mydatabase.Database) *Service {
	return &Service{
		repo:        repo,
		ActiveUsers: make(map[string]ActiveUser),
		mutex:       sync.Mutex{},
		activeMu:    sync.Mutex{},
	}
}
