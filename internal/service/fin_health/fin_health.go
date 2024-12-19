package fin_health

import (
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
)

type Service struct {
	repo *mydb.Database
}

func NewService(repo *mydb.Database) *Service {
	return &Service{
		repo: repo,
	}
}

type Health interface {
	ExpenditureDelta(userID string) (float64, error)
}
