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
