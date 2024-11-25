package models

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

type WelfareFund int

const (
	Planned WelfareFund = iota
	Unplanned
)
