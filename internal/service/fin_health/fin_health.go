package fin_health

import (
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
)

const (
	liquid   = true
	illiquid = false
)

const (
	saving     = "saving"
	investment = "investment"
	loan       = "loan"
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
	ExpensePropensity(userID string) (float64, error)
	LiquidFundRatio(userID string) (float64, error)
	IlliquidFundRatio(userID string) (float64, error)
	SavingsToIncomeRatio(userID string) (float64, error)
	SavingDelta(userID string) (float64, error)
}
