// Package repository provides basic financial repository functionality.
package repository

import (
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"time"
)

type WealthFundModel struct {
	DB *mydb.Database
}

func (m *WealthFundModel) Create(wealthFund *models.WealthFund) (int64, error) {
	parsedDate, err := time.Parse("2006-01-02", wealthFund.Date)
	if err != nil {
		return 0, err
	}

	var wealthFundID int64
	err1 := m.DB.QueryRow("INSERT INTO wealth_fund (amount, date, planned, user_id, currency_code, connected_account, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		wealthFund.Amount, parsedDate, wealthFund.PlannedStatus, wealthFund.UserID, wealthFund.Currency, wealthFund.ConnectedAccount, wealthFund.CategoryID).Scan(&wealthFundID)
	if err1 != nil {
		return 0, err1
	}
	return wealthFundID, nil
}

func (m *WealthFundModel) Delete(id, userID string) error {
	result, err := m.DB.Exec("DELETE FROM wealth_fund WHERE id = $1 AND user_id = $2", id, userID)
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
		return fmt.Errorf("%w: no wealth fund found with id %s for user %s", myerrors.ErrNotFound, id, userID)
	}

	return nil
}

func (m *WealthFundModel) Update(wealthFund *models.WealthFund) error {
	q := `
		UPDATE wealth_fund SET 
		   amount=$1, 
		   date=$2, 
		   planned=$3, 
		   currency_code=$4,
		   connected_account=$5, 
           category_id=$6
	   WHERE id=$7 AND user_id=$8`

	result, err := m.DB.Exec(q, wealthFund.Amount, wealthFund.Date, wealthFund.PlannedStatus, wealthFund.Currency,
		wealthFund.ConnectedAccount, wealthFund.CategoryID, wealthFund.ID, wealthFund.UserID)

	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: no wealth func found with id %s for user %s", myerrors.ErrNotFound)
	}

	return nil
}

func (m *WealthFundModel) ListByUserID(userID string) ([]models.WealthFund, error) {
	rows, err := m.DB.Query("SELECT id, amount, date, planned, currency_code, connected_account FROM wealth_fund WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wealthFunds []models.WealthFund
	for rows.Next() {
		var wealthFund models.WealthFund
		if err := rows.Scan(&wealthFund.ID, &wealthFund.Amount, &wealthFund.Date); err != nil {
			return nil, err
		}
		wealthFund.UserID = userID
		wealthFunds = append(wealthFunds, wealthFund)
	}

	return wealthFunds, nil
}
