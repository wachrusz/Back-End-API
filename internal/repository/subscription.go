//go:build !exclude_swagger
// +build !exclude_swagger

// Package repository provides basic financial repository functionality.
package repository

import (
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
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

func CreateSubscription(subscription *Subscription) (int64, error) {
	parsedDateStart, err := time.Parse("2006-01-02", subscription.StartDate)
	if err != nil {
		log.Println("Error parsing date:", err)
		return 0, err
	}

	parsedDateEnd, err1 := time.Parse("2006-01-02", subscription.EndDate)
	if err1 != nil {
		log.Println("Error parsing date:", err)
		return 0, err1
	}

	var subscriptionID int64
	err2 := mydb.GlobalDB.QueryRow("INSERT INTO subscriptions (user_id, start_date, end_date, is_active) VALUES ($1, $2, $3, $4) RETURNING id",
		subscription.UserID, parsedDateStart, parsedDateEnd, subscription.IsActive).Scan(&subscriptionID)
	if err2 != nil {
		log.Println("Error creating income:", err1)
		return 0, err2
	}
	return subscriptionID, nil
}
