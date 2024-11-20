package repository

import (
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"time"
)

// Models is a facade for all models interfaces.
type Models struct {
	Accounts Accounts
	Expenses Expenses
}

func New(db *mydb.Database) *Models {
	return &Models{
		Accounts: &AccountModel{db},
		Expenses: &ExpenseModel{db},
	}
}

type Accounts interface {
	Create(account *models.ConnectedAccount) (int64, error)
	Delete(id, userID string) error
	Update(editedAccount *models.ConnectedAccount) error
}

type Expenses interface {
	Create(expense *models.Expense) (int64, error)
	GetByUserID(userID string) ([]models.Expense, error)
	GetForMonth(userID string, month time.Month, year int) (float64, float64, error)
	GetMonthlyIncrease(userID string) (int, int, error)
	Delete(id, userID string) error
	Update(editedExpense *models.Expense) error
}
