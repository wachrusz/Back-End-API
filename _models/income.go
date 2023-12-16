//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "backEndAPI/_mydatabase"
	"database/sql"

	"log"
	"time"
)

type Income struct {
	ID         string  `json:"id"`
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
	Planned    bool    `json:"planned"`
	UserID     string  `json:"user_id"`
	CategoryID string  `json:"category_id"`
}

func CreateIncome(income *Income) error {
	parsedDate, err := time.Parse("2006-01-02", income.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return err
	}
	_, err1 := mydb.GlobalDB.Exec("INSERT INTO income (amount, date, planned, user_id, category) VALUES ($1, $2, $3, $4, $5)",
		income.Amount, parsedDate, income.Planned, income.UserID, income.CategoryID)
	if err1 != nil {
		log.Println("Error creating income:", err1)
		return err1
	}
	_, err = mydb.GlobalDB.Exec("INSERT INTO operations (user_id, description, amount, date, category, operation_type) VALUES ($1, $2, $3, $4, $5, $6)",
		income.UserID, "Доход", income.Amount, parsedDate, income.CategoryID, income.CategoryID)
	if err != nil {
		log.Println("Error creating income operation:", err)
		return err
	}
	return nil
}

func GetIncomesByUserID(userID string) ([]Income, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, amount, date, planned, category FROM income WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying incomes:", err)
		return nil, err
	}
	defer rows.Close()

	var incomes []Income
	for rows.Next() {
		var income Income
		if err := rows.Scan(&income.ID, &income.Amount, &income.Date, &income.Planned, &income.CategoryID); err != nil {
			log.Println("Error scanning income row:", err)
			return nil, err
		}
		income.UserID = userID
		incomes = append(incomes, income)
	}

	return incomes, nil
}

func GetIncomeForMonth(userID string, month time.Month, year int) (float64, float64, error) {
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
	err := mydb.GlobalDB.QueryRow(query, userID, int(month), year).Scan(&totalIncome, &plannedIncome)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting income for month: %v", err)
		return 0, 0, err
	}

	return totalIncome, plannedIncome, nil
}

func GetMonthlyIncomeIncrease(userID string) (int, int, error) {
	currentDate := time.Now()

	currentMonth := currentDate.Month()
	currentYear := currentDate.Year()

	previousMonth := currentMonth - 1
	previousYear := currentYear

	if currentMonth == time.January {
		previousMonth = time.December
		previousYear--
	}

	currentMonthIncome, currentMonthPlanned, err := GetIncomeForMonth(userID, currentMonth, currentYear)
	if err != nil {
		log.Printf("Error fetching current month income: %v", err)
		return 0, 0, err
	}

	previousMonthIncome, _, err := GetIncomeForMonth(userID, previousMonth, previousYear)
	if err != nil {
		log.Printf("Error fetching previous month income: %v", err)
		return 0, 0, err
	}

	return int(((currentMonthIncome / previousMonthIncome) - 1) * 100), int(((currentMonthPlanned / currentMonthIncome) - 1) * 100), nil
}
