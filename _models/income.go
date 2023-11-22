package models

import (
	mydb "backEndAPI/_mydatabase"

	"log"
	"time"
)

type Income struct {
	ID      string  `json:"id"`
	Amount  float64 `json:"amount"`
	Date    string  `json:"date"`
	Planned bool    `json:"planned"`
	UserID  string  `json:"user_id"`
}

func CreateIncome(income *Income) error {
	parsedDate, err := time.Parse("2006-01-02", income.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return err
	}
	_, err1 := mydb.GlobalDB.Exec("INSERT INTO income (amount, date, planned, user_id) VALUES ($1, $2, $3, $4)",
		income.Amount, parsedDate, income.Planned, income.UserID)
	if err1 != nil {
		log.Println("Error creating income:", err1)
		return err1
	}
	return nil
}

// GetIncomesByUserID возвращает список записей о доходе для определенного пользователя.
func GetIncomesByUserID(userID string) ([]Income, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, amount, date, planned FROM income WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying incomes:", err)
		return nil, err
	}
	defer rows.Close()

	var incomes []Income
	for rows.Next() {
		var income Income
		if err := rows.Scan(&income.ID, &income.Amount, &income.Date, &income.Planned); err != nil {
			log.Println("Error scanning income row:", err)
			return nil, err
		}
		income.UserID = userID
		incomes = append(incomes, income)
	}

	return incomes, nil
}
