package user

import (
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/pkg/mydatabase"
	"sync"
)

type UserService interface {
}

type Service struct {
	email *email.EmailServce
	repo  *mydatabase.Database
	mutex sync.Mutex
}

func NewService(email *email.EmailServce, repo *mydatabase.Database) *Service {
	return &Service{
		email: email,
		repo:  repo,
		mutex: sync.Mutex{},
	}
}
