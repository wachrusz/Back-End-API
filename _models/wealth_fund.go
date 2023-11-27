//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "backEndAPI/_mydatabase"
	"log"
	"time"
)

type WealthFund struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Date   string  `json:"date"`
	UserID string  `json:"user_id"`
}

// @Summary Create wealth fund entry
// @Description Create a new wealth fund entry.
// @Tags WealthFund
// @Accept json
// @Produce json
// @Param wealthFund body WealthFund true "Wealth fund details"
// @Success 201 {string} string "Wealth fund created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating wealth fund"
// @Router /models/wealth-fund [post]
func CreateWealthFund(wealthFund *WealthFund) error {
	parsedDate, err := time.Parse("2006-01-02", wealthFund.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return err
	}

	_, err1 := mydb.GlobalDB.Exec("INSERT INTO wealth_fund (amount, date, user_id) VALUES ($1, $2, $3)",
		wealthFund.Amount, parsedDate, wealthFund.UserID)
	if err1 != nil {
		log.Println("Error creating wealthFund:", err)
		return err1
	}
	return nil
}

// @Summary Get wealth funds by user ID
// @Description Get a list of wealth funds for a specific user.
// @Tags WealthFund
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} WealthFund "List of wealth funds"
// @Failure 500 {string} string "Error querying wealth funds"
// @Router /models/wealth-fund/{userID} [get]
func GetWealthFundsByUserID(userID string) ([]WealthFund, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, amount, date FROM wealthFund WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying wealthFunds:", err)
		return nil, err
	}
	defer rows.Close()

	var wealthFunds []WealthFund
	for rows.Next() {
		var wealthFund WealthFund
		if err := rows.Scan(&wealthFund.ID, &wealthFund.Amount, &wealthFund.Date); err != nil {
			log.Println("Error scanning wealthFund row:", err)
			return nil, err
		}
		wealthFund.UserID = userID
		wealthFunds = append(wealthFunds, wealthFund)
	}

	return wealthFunds, nil
}
