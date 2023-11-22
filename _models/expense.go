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

// GetExpensesByUserID возвращает список записей о доходе для определенного пользователя.
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
