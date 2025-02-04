// Package repository provides basic financial repository functionality.
package repository

import (
	"database/sql"
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"log"
	"time"
)

type IncomeModel struct {
	DB *mydb.Database
}

func (m *IncomeModel) Create(income *models.Income) (int64, error) {
	parsedDate, err := time.Parse("2006-01-02", income.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return 0, err
	}

	var incomeID int64
	err = m.DB.QueryRow("INSERT INTO income (amount, date, planned, user_id, category, sender, connected_account, currency_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		income.Amount, parsedDate, income.Planned, income.UserID, income.CategoryID, income.Sender, income.BankAccount, income.Currency).Scan(&incomeID)
	if err != nil {
		return 0, err
	}
	_, err = m.DB.Exec("INSERT INTO operations (user_id, description, amount, date, category, operation_type) VALUES ($1, $2, $3, $4, $5, $6)",
		income.UserID, "Доход", income.Amount, parsedDate, income.CategoryID, income.CategoryID)
	if err != nil {
		return 0, err
	}
	return incomeID, nil
}

func (m *IncomeModel) ListByUserID(userID string) ([]models.Income, error) {
	rows, err := m.DB.Query("SELECT id, amount, date, planned, category, sender, connected_account, currency_code FROM income WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []models.Income
	for rows.Next() {
		var income models.Income
		if err := rows.Scan(&income.ID, &income.Amount, &income.Date, &income.Planned, &income.CategoryID, &income.Sender, &income.BankAccount, &income.Currency); err != nil {
			return nil, err
		}
		income.UserID = userID
		incomes = append(incomes, income)
	}

	return incomes, nil
}

func (m *IncomeModel) ListByMonth(userID string, month time.Month, year int) (float64, float64, error) {
	query := `
		SELECT
			COALESCE(SUM(amount), 0) AS total_income,
			COALESCE(SUM(CASE WHEN planned THEN amount ELSE 0 END), 0) AS planned_income
		FROM income
		WHERE user_id = $1
		AND EXTRACT(MONTH FROM date) = $2
		AND EXTRACT(YEAR FROM date) = $3
	`

	var totalIncome, plannedIncome float64
	err := m.DB.QueryRow(query, userID, int(month), year).Scan(&totalIncome, &plannedIncome)
	if err != nil && err != sql.ErrNoRows {
		return 0, 0, err
	}

	return totalIncome, plannedIncome, nil
}

func (m *IncomeModel) GetMonthlyIncomeIncrease(userID string) (int, int, error) {
	currentDate := time.Now()

	currentMonth := currentDate.Month()
	currentYear := currentDate.Year()

	previousMonth := currentMonth - 1
	previousYear := currentYear

	if currentMonth == time.January {
		previousMonth = time.December
		previousYear--
	}

	currentMonthIncome, currentMonthPlanned, err := m.ListByMonth(userID, currentMonth, currentYear)
	if err != nil {
		return 0, 0, err
	}

	previousMonthIncome, _, err := m.ListByMonth(userID, previousMonth, previousYear)
	if err != nil {
		return 0, 0, err
	}

	return int(((currentMonthIncome / previousMonthIncome) - 1) * 100), int(((currentMonthPlanned / currentMonthIncome) - 1) * 100), nil
}

func (m *IncomeModel) Delete(id, userID string) error {
	result, err := m.DB.Exec("DELETE FROM income WHERE id = $1 AND user_id = $2", id, userID)
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
		return fmt.Errorf("%w: no income found with id %s for user %s", myerrors.ErrNotFound, id, userID)
	}

	return nil
}

func (m *IncomeModel) Update(income *models.Income) error {
	q := `
		UPDATE income SET 
			amount = $1, 
			date = $2, 
			planned = $3, 
			category = $4, 
			sender = $5, 
			connected_account = $6, 
			currency_code = $7
		WHERE id = $8 AND user_id = $9`

	result, err := m.DB.Exec(q, income.Amount, income.Date, income.Planned, income.CategoryID,
		income.Sender, income.BankAccount, income.Currency, income.ID, income.UserID)

	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: no income found with id %s for user %s", myerrors.ErrNotFound, income.ID, income.UserID)
	}

	return nil
}
