//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	"log"
	"time"

	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
)

type Goal struct {
	ID           string    `json:"id"`
	Goal         string    `json:"goal"`
	Need         float64   `json:"need"`
	CurrentState float64   `json:"current_state"`
	Currency     string    `json:"currency"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	UserID       string    `json:"user_id"`
}

func CreateGoal(goal *Goal) (int64, error) {
	var goalID int64
	err := mydb.GlobalDB.QueryRow("INSERT INTO goal (goal, need, current_state, end_date, currency, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		goal.Goal, goal.Need, goal.CurrentState, goal.EndDate, goal.Currency, goal.UserID).Scan(&goalID)
	if err != nil {
		log.Println("Error creating goal:", err)
		return 0, err
	}
	return goalID, nil
}

func UpdateGoal(goal *Goal) (int64, error) {
	var goalID int64
	err := mydb.GlobalDB.QueryRow(`
		UPDATE goal 
		SET 
			goal = $1,
			need = $2,
			end_date = $3
		WHERE id = $4
		RETURNING id;
	`, goal.Goal, goal.Need, goal.EndDate, goal.ID).Scan(&goalID)

	if err != nil {
		log.Println("Error updating goal:", err)
		return 0, err
	}
	return goalID, nil
}
func GetGoalsByUserID(userID string) ([]Goal, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, goal, need, current_state, start_date, end_date, currency FROM goal WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying goals:", err)
		return nil, err
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var goal Goal
		if err := rows.Scan(&goal.ID, &goal.Goal, &goal.Need, &goal.CurrentState, goal.StartDate, &goal.EndDate, &goal.Currency); err != nil {
			log.Println("Error scanning goal row:", err)
			return nil, err
		}
		goal.UserID = userID
		goals = append(goals, goal)
	}

	return goals, nil
}
