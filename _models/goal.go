//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "backEndAPI/_mydatabase"
	"log"
)

type Goal struct {
	ID           string  `json:"id"`
	Goal         string  `json:"goal"`
	Need         float64 `json:"need"`
	CurrentState float64 `json:"current_state"`
	UserID       string  `json:"user_id"`
}

func CreateGoal(goal *Goal) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO goal (goal, need, current_state, user_id) VALUES ($1, $2, $3, $4)",
		goal.Goal, goal.Need, goal.CurrentState, goal.UserID)
	if err != nil {
		log.Println("Error creating goal:", err)
		return err
	}
	return nil
}

func GetGoalsByUserID(userID string) ([]Goal, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, goal, need, current_state FROM goal WHERE user_id = $1", userID)
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
