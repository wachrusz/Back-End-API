package repository

import (
	"fmt"

	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
)

type AccountModel struct {
	DB *mydb.Database
}

func (m *AccountModel) Create(account *models.ConnectedAccount) (int64, error) {
	var connectedAccountID int64
	err := m.DB.QueryRow(
		`INSERT INTO connected_accounts 
		(user_id, bank_id, account_number, account_type, state, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`,
		account.UserID, account.BankID, account.AccountNumber, account.AccountType, account.AccountState).Scan(&connectedAccountID)
	if err != nil {
		return 0, err
	}
	return connectedAccountID, nil

}

func (m *AccountModel) Delete(id, userID string) error {
	result, err := m.DB.Exec("DELETE FROM connected_accounts WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		// Возвращаем обернутую ошибку, если запрос завершился с ошибкой
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Возвращаем обернутую ошибку при невозможности получить количество строк
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	if rowsAffected == 0 {
		// Возвращаем ошибку, если запись не найдена или не принадлежит пользователю
		return fmt.Errorf("%w: no account found with id %s for user %s", myerrors.ErrNotFound, id, userID)
	}

	return nil
}

func (m *AccountModel) Update(account *models.ConnectedAccount) error {
	result, err := m.DB.Exec("UPDATE connected_accounts SET bank_id=$1, account_number=$2, account_type=$3, updated_at=NOW() WHERE id = $4 AND user_id = $5",
		account.BankID, account.AccountNumber, account.AccountType, account.ID, account.UserID)

	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: no account found with id %s for user %s", myerrors.ErrNotFound, account.ID, account.UserID)
	}

	return nil
}
