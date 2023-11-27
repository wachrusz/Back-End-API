//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "backEndAPI/_mydatabase"

	"log"
	"time"
)

type Expense struct {
	ID      string  `json:"id"`
	Amount  float64 `json:"amount"`
	Date    string  `json:"date"`
	Planned bool    `json:"planned"`
	UserID  string  `json:"user_id"`
}

// @Summary Create expense entry
// @Description Create a new expense entry.
// @Tags Expense
// @Accept json
// @Produce json
// @Param expense body Expense true "Expense details"
// @Success 201 {string} string "Expense created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating expense"
// @Router /analytics/expence [post]
func CreateExpense(expense *Expense) error {
	parsedDate, err := time.Parse("2006-01-02", expense.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return err
	}
	_, err1 := mydb.GlobalDB.Exec("INSERT INTO expense (amount, date, planned, user_id) VALUES ($1, $2, $3, $4)",
		expense.Amount, parsedDate, expense.Planned, expense.UserID)
	if err1 != nil {
		log.Println("Error creating expense:", err)
		return err1
	}
	return nil
}

// @Summary Get expenses by user ID
// @Description Get a list of expenses for a specific user.
// @Tags Expense
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} Expense "List of expenses"
// @Failure 500 {string} string "Error querying expenses"
// @Router /analytics/expence/{userID} [get]
func GetExpensesByUserID(userID string) ([]Expense, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, amount, date, planned FROM expense WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying expenses:", err)
		return nil, err
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(&expense.ID, &expense.Amount, &expense.Date, &expense.Planned); err != nil {
			log.Println("Error scanning expense row:", err)
			return nil, err
		}
		expense.UserID = userID
		expenses = append(expenses, expense)
	}

	return expenses, nil
}
