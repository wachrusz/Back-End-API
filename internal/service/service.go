package service

import (
	"github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
	"github.com/wachrusz/Back-End-API/internal/service/currency"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/token"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	"github.com/wachrusz/Back-End-API/pkg/rabbit"
)

type Services struct {
	Users      user.Users
	Categories categories.Categories
	Emails     email.Emails
	Currency   currency.CurrencyService
	Tokens     token.Tokens
}

type Dependencies struct {
	Repo                  *mydatabase.Database
	Mailer                rabbit.Mailer
	AccessTokenDurMinutes int
}

func NewServices(deps Dependencies) (*Services, error) {
	cur, err := currency.NewService(deps.Repo)
	if err != nil {
		return nil, err
	}
	e := email.NewService(deps.Repo, deps.Mailer)
	cat := categories.NewService(deps.Repo, cur)
	u := user.NewService(deps.Repo, cat)
	t := token.NewService(deps.Repo, e, u, deps.AccessTokenDurMinutes)
	return &Services{
		Users:      u,
		Categories: cat,
		Emails:     e,
		//Reports:    report.NewService(deps.Repo),
		Currency: cur,
		Tokens:   t,
	}, nil
}
