//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	"log"

	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
)

type Goal struct {
	ID           string  `json:"id"`
	Goal         string  `json:"goal"`
	Need         float64 `json:"need"`
	CurrentState float64 `json:"current_state"`
	Currency     string  `json:"currency"`
	UserID       string  `json:"user_id"`
}

func CreateGoal(goal *Goal) (int64, error) {
	var goalID int64
	err := mydb.GlobalDB.QueryRow("INSERT INTO goal (goal, need, current_state, currency, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		goal.Goal, goal.Need, goal.CurrentState, goal.Currency, goal.UserID).Scan(&goalID)
	if err != nil {
		log.Println("Error creating goal:", err)
		return 0, err
	}
	return goalID, nil
}

func GetGoalsByUserID(userID string) ([]Goal, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, goal, need, current_state, currency FROM goal WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying goals:", err)
		return nil, err
	}
	defer rows.Close()

	var goals []Goal
	for rows.Next() {
		var goal Goal
		if err := rows.Scan(&goal.ID, &goal.Need, &goal.CurrentState, &goal.Goal); err != nil {
			log.Println("Error scanning goal row:", err)
			return nil, err
		}
		goal.UserID = userID
		goals = append(goals, goal)
	}

	return goals, nil
}
