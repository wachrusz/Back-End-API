//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"log"
	"time"
)

type WelfareFund int

const (
	Planned WelfareFund = iota
	Unplanned
)

type WealthFund struct {
	ID               string      `json:"id"`
	Amount           float64     `json:"amount"`
	Date             string      `json:"date"`
	PlannedStatus    WelfareFund `json:"planned"`
	Currency         string      `json:"currency"`
	ConnectedAccount string      `json:"bank_account"`
	CategoryID       string      `json:"category_id"`
	UserID           string      `json:"user_id"`
}

func CreateWealthFund(wealthFund *WealthFund) (int64, error) {
	parsedDate, err := time.Parse("2006-01-02", wealthFund.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return 0, err
	}

	var wealthFundID int64
	err1 := mydb.GlobalDB.QueryRow("INSERT INTO wealth_fund (amount, date, planned, user_id, currency_code, connected_account, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		wealthFund.Amount, parsedDate, wealthFund.PlannedStatus, wealthFund.UserID, wealthFund.Currency, wealthFund.ConnectedAccount, wealthFund.CategoryID).Scan(&wealthFundID)
	if err1 != nil {
		log.Println("Error creating wealthFund:", err)
		return 0, err1
	}
	return wealthFundID, nil
}

func GetWealthFundsByUserID(userID string) ([]WealthFund, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, amount, date, planned, currency_code, connected_account FROM wealth_fund WHERE user_id = $1", userID)
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
