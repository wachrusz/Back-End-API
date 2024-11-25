package models

type Income struct {
	ID          string  `json:"id"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Planned     bool    `json:"planned"`
	UserID      string  `json:"user_id"`
	CategoryID  string  `json:"category_id"`
	Sender      string  `json:"sender"`
	BankAccount string  `json:"bank_account"`
	Currency    string  `json:"currency"`
}
