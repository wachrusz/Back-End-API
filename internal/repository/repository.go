package repository

import (
	"time"

	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
)

// Models is a facade for all models interfaces.
type Models struct {
	Accounts    Accounts
	Expenses    Expenses
	Goals       Goals
	Incomes     Incomes
	WealthFunds WealthFunds
}

func New(db *mydb.Database) *Models {
	return &Models{
		Accounts:    &AccountModel{db},
		Expenses:    &ExpenseModel{db},
		Goals:       &GoalModel{db},
		Incomes:     &IncomeModel{db},
		WealthFunds: &WealthFundModel{db},
	}
}

type Accounts interface {
	Create(account *models.ConnectedAccount) (int64, error)
	Update(account *models.ConnectedAccount) error
	Delete(id, userID string) error
}

type Expenses interface {
	Create(expense *models.Expense) (int64, error)
	Update(expense *models.Expense) error
	Delete(id, userID string) error
	GetByUserID(userID string) ([]models.Expense, error)
	GetForMonth(userID string, month time.Month, year int) (float64, float64, error)
	GetMonthlyIncrease(userID string) (int, int, error)
}

type Goals interface {
	Create(goal *models.Goal) (int64, error)
	Update(goal *models.Goal) error
	Delete(id, userID string) error
	GetByUserID(userID string) ([]models.Goal, error)
}

type Incomes interface {
	Create(income *models.Income) (int64, error)
	Update(income *models.Income) error
	Delete(id, userID string) error
	GetIncomesByUserID(userID string) ([]models.Income, error)
	GetIncomeForMonth(userID string, month time.Month, year int) (float64, float64, error)
	GetMonthlyIncomeIncrease(userID string) (int, int, error)
}

type WealthFunds interface {
	Create(wealthFund *models.WealthFund) (int64, error)
	Update(wealthFund *models.WealthFund) error
	Delete(id, userID string) error
	GetByUserID(userID string) ([]models.WealthFund, error)
}
