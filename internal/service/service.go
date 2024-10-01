package service

import (
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	"github.com/wachrusz/Back-End-API/pkg/mydatabase"
)

type Services struct {
	Users  *user.Service
	Emails *email.EmailServce
}

type Dependencies struct {
	Repo *mydatabase.Database
}

func NewServices(deps Dependencies) *Services {
	mailer := email.NewService()
	return &Services{
		Users:  user.NewService(mailer, deps.Repo),
		Emails: mailer,
	}
}
