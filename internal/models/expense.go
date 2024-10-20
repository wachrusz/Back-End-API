//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	"database/sql"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"log"
	"time"
)

type Expense struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Planned     bool    `json:"planned"`
	UserID      string  `json:"user_id"`
	CategoryID  string  `json:"category_id"`
	SentTo      string  `json:"sent_to"`
	BankAccount string  `json:"bank_account"`
	Currency    string  `json:"currency"`
}

func CreateExpense(expense *Expense) (int64, error) {
	parsedDate, err := time.Parse("2006-01-02", expense.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return 0, err
	}

	var expenseID int64
	err = mydb.GlobalDB.QueryRow("INSERT INTO expense (amount, date, planned, user_id, category, sent_to, connected_account, currency_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		expense.Amount, parsedDate, expense.Planned, expense.UserID, expense.CategoryID, expense.SentTo, expense.BankAccount, expense.Currency).Scan(&expenseID)
	if err != nil {
		log.Println("Error creating expense:", err)
		return 0, err
	}
	_, err = mydb.GlobalDB.Exec("INSERT INTO operations (user_id, description, amount, date, category, operation_type) VALUES ($1, $2, $3, $4, $5, $6)",
		expense.UserID, "Расход", expense.Amount, parsedDate, expense.CategoryID, expense.CategoryID)
	if err != nil {
		log.Println("Error creating expense operation:", err)
		return 0, err
	}
	return expenseID, nil
}

func GetExpensesByUserID(userID string) ([]Expense, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, amount, date, planned, category, sent_to, connected_account, currency_code FROM expense WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying expenses:", err)
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(&expense.ID, &expense.Amount, &expense.Date, &expense.Planned, &expense.CategoryID, &expense.SentTo, &expense.BankAccount, &expense.Currency); err != nil {
			log.Println("Error scanning expense row:", err)
			return nil, err
		}
		expense.UserID = userID
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func GetExpenseForMonth(userID string, month time.Month, year int) (float64, float64, error) {

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
	err := mydb.GlobalDB.QueryRow(query, userID, int(month), year).Scan(&totalExpense, &plannedExpense)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error getting expense for month: %v", err)
		return 0, 0, err
	}

	return totalExpense, plannedExpense, nil
}

func GetMonthlyExpenseIncrease(userID string) (int, int, error) {
	currentDate := time.Now()

	currentMonth := currentDate.Month()
	currentYear := currentDate.Year()

	previousMonth := currentMonth - 1
	previousYear := currentYear

	if currentMonth == time.January {
		previousMonth = time.December
		previousYear--
	}

	currentMonthExpense, currentMonthPlanned, err := GetExpenseForMonth(userID, currentMonth, currentYear)
	if err != nil {
		log.Printf("Error fetching current month income: %v", err)
		return 0, 0, err
	}

	previousMonthExpense, _, err := GetExpenseForMonth(userID, previousMonth, previousYear)
	if err != nil {
		log.Printf("Error fetching previous month income: %v", err)
		return 0, 0, err
	}

	return int(((currentMonthExpense / previousMonthExpense) - 1) * 100), int(((currentMonthPlanned / previousMonthExpense) - 1) * 100), nil
}
