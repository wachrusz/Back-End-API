package models

import "time"

type Goal struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	Date             time.Time `json:"date"`
	Months           int       `json:"months"`
	UserID           int64     `json:"user_id"`
	AdditionalMonths int       `json:"additional_months"`
	IsCompleted      bool      `json:"is_completed"`
}

type GoalTransaction struct {
	ID          int64     `json:"id"`
	GoalID      int64     `json:"goal_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Date        time.Time `json:"date"`
	Planned     bool      `json:"planned"`
	BankAccount string    `json:"bank_account"`
	//Type        string    `json:"type"`
}

type GoalDetails struct {
	Goal           Goal    `json:"goal"`
	Month          int     `json:"month"`
	MonthlyPayment float64 `json:"monthly_payment"`
	CurrentPayment float64 `json:"current_payment"`
	CurrentNeed    float64 `json:"current_need"`
	Gathered       float64 `json:"gathered"`
}

type GoalTrackerInfo struct {
	Goal         *Goal              `json:"goal"`
	Transactions []*GoalTransaction `json:"transactions"`
}
