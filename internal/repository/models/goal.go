package models

import "time"

type Goal struct {
	ID           int64     `json:"id"`
	Goal         string    `json:"goal"`
	Need         float64   `json:"need"`
	CurrentState float64   `json:"current_state"`
	Currency     string    `json:"currency"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	UserID       string    `json:"user_id"`
}
