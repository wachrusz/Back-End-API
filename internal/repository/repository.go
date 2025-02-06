package repository

import (
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"time"

	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
)

// Models is a facade for all models interfaces.
type Models struct {
	Accounts          AccountRepo
	Expenses          ExpenseRepo
	Goals             GoalRepo
	GoalsTransactions GoalTransactionRepo
	Incomes           IncomeRepo
	WealthFunds       WealthFundRepo
	Subscriptions     SubscriptionRepo
}

func New(db *mydb.Database) *Models {
	return &Models{
		Accounts:          &AccountModel{db},
		Expenses:          &ExpenseModel{db},
		Goals:             &GoalModel{db},
		GoalsTransactions: &GoalTransactionModel{db},
		Incomes:           &IncomeModel{db},
		WealthFunds:       &WealthFundModel{db},
		Subscriptions:     &SubscriptionModel{db},
	}
}

type AccountRepo interface {
	Create(account *models.ConnectedAccount) (int64, error)
	Update(account *models.ConnectedAccount) error
	Delete(id, userID string) error
}

type ExpenseRepo interface {
	Create(expense *models.Expense) (int64, error)
	Update(expense *models.Expense) error
	Delete(id, userID string) error
	ListByUserID(userID string) ([]models.Expense, error)
	GetForMonth(userID string, month time.Month, year int) (float64, float64, error)
	GetMonthlyIncrease(userID string) (int, int, error)
}

type GoalRepo interface {
	Create(goal *models.Goal) (id int64, err error)
	Update(goal *models.Goal) error
	Delete(id int64, userID int64) error
	ListByUserID(userID int64) ([]models.Goal, error)
	Details(id int64, userID int64) (*models.GoalDetails, error)
	TrackerInfo(userID int64, limitStr, offsetStr int) ([]*models.GoalTrackerInfo, *jsonresponse.Metadata, error)
}

type GoalTransactionRepo interface {
	Create(transaction *models.GoalTransaction, userID int64) (id int64, err error)
}

type IncomeRepo interface {
	Create(income *models.Income) (int64, error)
	Update(income *models.Income) error
	Delete(id, userID string) error
	ListByUserID(userID string) ([]models.Income, error)
	ListByMonth(userID string, month time.Month, year int) (float64, float64, error)
	GetMonthlyIncomeIncrease(userID string) (int, int, error)
}

type WealthFundRepo interface {
	Create(wealthFund *models.WealthFund) (int64, error)
	Update(wealthFund *models.WealthFund) error
	Delete(id, userID string) error
	ListByUserID(userID string) ([]models.WealthFund, error)
}

type SubscriptionRepo interface {
	Create(subscription *models.Subscription) (int64, error)
	Update(subscription *models.Subscription) error
	Delete(id, userID string) error
}
