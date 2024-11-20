// Package repository provides basic financial repository functionality.
package repository

import (
	"database/sql"
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"time"
)

type ExpenseModel struct {
	DB *mydb.Database
}

func (m *ExpenseModel) Create(expense *models.Expense) (int64, error) {
	parsedDate, err := time.Parse("2006-01-02", expense.Date)
	if err != nil {
		return 0, err
	}

	var expenseID int64
	err = m.DB.QueryRow("INSERT INTO expense (amount, date, planned, user_id, category, sent_to, connected_account, currency_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		expense.Amount, parsedDate, expense.Planned, expense.UserID, expense.CategoryID, expense.SentTo, expense.BankAccount, expense.Currency).Scan(&expenseID)

	if err != nil {
		return 0, err
	}

	_, err = m.DB.Exec("INSERT INTO operations (user_id, description, amount, date, category, operation_type) VALUES ($1, $2, $3, $4, $5, $6)",
		expense.UserID, "Расход", expense.Amount, parsedDate, expense.CategoryID, expense.CategoryID)

	if err != nil {
		return 0, err
	}
	return expenseID, nil
}

func (m *ExpenseModel) GetByUserID(userID string) ([]models.Expense, error) {
	rows, err := m.DB.Query("SELECT id, amount, date, planned, category, sent_to, connected_account, currency_code FROM expense WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		if err := rows.Scan(&expense.ID, &expense.Amount, &expense.Date, &expense.Planned, &expense.CategoryID, &expense.SentTo, &expense.BankAccount, &expense.Currency); err != nil {
			return nil, err
		}
		expense.UserID = userID
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func (m *ExpenseModel) GetForMonth(userID string, month time.Month, year int) (float64, float64, error) {
	query := `
		SELECT
			COALESCE(SUM(amount), 0) AS total_expense,
			COALESCE(SUM(CASE WHEN planned THEN amount ELSE 0 END), 0) AS planned_expense
		FROM expense
		WHERE user_id = $1
		AND EXTRACT(MONTH FROM date) = $2
		AND EXTRACT(YEAR FROM date) = $3
	`

	var totalExpense, plannedExpense float64
	err := m.DB.QueryRow(query, userID, int(month), year).Scan(&totalExpense, &plannedExpense)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, err
	}

	return totalExpense, plannedExpense, nil
}

func (m *ExpenseModel) GetMonthlyIncrease(userID string) (int, int, error) {
	currentDate := time.Now()

	currentMonth := currentDate.Month()
	currentYear := currentDate.Year()

	previousMonth := currentMonth - 1
	previousYear := currentYear

	if currentMonth == time.January {
		previousMonth = time.December
		previousYear--
	}

	currentMonthExpense, currentMonthPlanned, err := m.GetForMonth(userID, currentMonth, currentYear)
	if err != nil {
		return 0, 0, err
	}

	previousMonthExpense, _, err := m.GetForMonth(userID, previousMonth, previousYear)
	if err != nil {
		return 0, 0, err
	}

	return int(((currentMonthExpense / previousMonthExpense) - 1) * 100), int(((currentMonthPlanned / previousMonthExpense) - 1) * 100), nil
}

func (m *ExpenseModel) Delete(id, userID string) error {
	result, err := m.DB.Exec("DELETE FROM expense WHERE id = $1 AND user_id = $2", id, userID)
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
		return fmt.Errorf("%w: no expense found with id %s for user %s", myerrors.ErrNotFound, id, userID)
	}

	return nil
}

func (m *ExpenseModel) Update(editedExpense *models.Expense) error {
	q := `
		UPDATE expense SET 
		   amount=$1, 
		   date=$2, 
		   planned=$3, 
		   category=$4, 
		   sent_to=$5, 
		   connected_account=$6, 
		   currency_code=$7 
	   WHERE id=$8 AND user_id=$9`

	result, err := m.DB.Exec(q, editedExpense.Amount, editedExpense.Date, editedExpense.Planned, editedExpense.CategoryID,
		editedExpense.SentTo, editedExpense.BankAccount, editedExpense.Currency, editedExpense.ID, editedExpense.UserID)

	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: no account found with id %s for user %s", myerrors.ErrNotFound, editedExpense.ID, editedExpense.UserID)
	}

	return nil
}
