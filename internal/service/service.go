package service

import (
	"github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
	"github.com/wachrusz/Back-End-API/internal/service/currency"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/token"
	"github.com/wachrusz/Back-End-API/internal/service/user"
)

type Services struct {
	Users      user.Users
	Categories categories.Categories
	Emails     email.Emails
	//Reports    *report.Service
	Currency currency.CurrencyService
	Tokens   token.Tokens
}

type Dependencies struct {
	Repo *mydatabase.Database
}

func NewServices(deps Dependencies) (*Services, error) {
	cur, err := currency.NewService(deps.Repo)
	if err != nil {
		return nil, err
	}
	u := user.NewService(deps.Repo)
	e := email.NewService(deps.Repo)
	cat := categories.NewService(deps.Repo, cur)
	t := token.NewService(deps.Repo, e, u)
	return &Services{
		Users:      u,
		Categories: cat,
		Emails:     e,
		//Reports:    report.NewService(deps.Repo),
		Currency: cur,
		Tokens:   t,
	}, nil
}
