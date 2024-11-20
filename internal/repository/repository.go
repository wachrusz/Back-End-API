package repository

import (
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
)

type Models struct {
	Accounts Accounts
}

func New(db *mydb.Database) *Models {
	return &Models{
		Accounts: &AccountModel{db},
	}
}

type Accounts interface {
	Create(account *models.ConnectedAccount) (int64, error)
	Delete(id, userID string) error
	Edit(editedAccount *models.ConnectedAccount) error
}
