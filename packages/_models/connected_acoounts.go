package models

import (
	"log"

	mydb "main/packages/_mydatabase"
)

// ConnectedAccount
type ConnectedAccount struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	BankID        string `json:"bank_id"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
}

func AddConnectedAccount(account *ConnectedAccount) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO connected_accounts (user_id, bank_id, account_number, account_type, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())",
		account.UserID, account.BankID, account.AccountNumber, account.AccountType)
	if err != nil {
		log.Println("Error adding connected_account:", err)
		return err
	}
	return nil

}

func DeleteConnectedAccount(userID string) error {
	_, err := mydb.GlobalDB.Exec("DELETE FROM connected_accounts WHERE user_id = $1",
		userID)
	if err != nil {
		log.Println("Error deleting connected_account:", err)
		return err
	}
	return nil
}
