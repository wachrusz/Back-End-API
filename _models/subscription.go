//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "backEndAPI/_mydatabase"
	"log"
	"time"
)

// Subscription представляет собой информацию о подписке пользователя.
type Subscription struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  bool   `json:"is_active"`
}

// @Summary Create subscription
// @Description Create a new subscription entry.
// @Tags Subscription
// @Accept json
// @Produce json
// @Param subscription body Subscription true "Subscription details"
// @Success 201 {string} string "Subscription created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating subscription"
// @Router /models/subscription [post]
func CreateSubscription(subscription *Subscription) error {
	parsedDateStart, err := time.Parse("2006-01-02", subscription.StartDate)
	if err != nil {
		log.Println("Error parsing date:", err)
		return err
	}

	parsedDateEnd, err1 := time.Parse("2006-01-02", subscription.EndDate)
	if err1 != nil {
		log.Println("Error parsing date:", err)
		return err1
	}

	_, err2 := mydb.GlobalDB.Exec("INSERT INTO subscriptions (user_id, start_date, end_date, is_active) VALUES ($1, $2, $3, $4)",
		subscription.UserID, parsedDateStart, parsedDateEnd, subscription.IsActive)
	if err2 != nil {
		log.Println("Error creating income:", err1)
		return err2
	}
	return nil
}
