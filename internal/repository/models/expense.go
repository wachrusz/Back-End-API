package models

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
