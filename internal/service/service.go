package service

import (
	"github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/report"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/user"
)

type Services struct {
	Users      *user.Service
	Categories *categories.Service
	Emails     *email.Service
	Reports    *report.Service
}

type Dependencies struct {
	Repo *mydatabase.Database
}

func NewServices(deps Dependencies) *Services {
	u := user.NewService(deps.Repo)
	return &Services{
		Users:      u,
		Categories: categories.NewService(deps.Repo),
		Emails:     email.NewService(deps.Repo, u),
		Reports:    report.NewService(deps.Repo),
	}
}
