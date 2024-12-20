package models

type ConnectedAccount struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	BankID          string  `json:"bank_id"`
	AccountName     string  `json:"account_name"`
	AccountNumber   string  `json:"account_number"`
	AccountType     string  `json:"account_type"`
	AccountState    float64 `json:"account_state"`
	AccountCurrency string  `json:"account_currency"`
}
