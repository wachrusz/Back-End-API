//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	"log"
	mydb "main/pkg/mydatabase"
)

type TrackingState struct {
	State  float64 `json:"state"`
	UserID string  `json:"user_id"`
}

func CreateTrackingState(trackingState *TrackingState) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO trackingState (state, user_id) VALUES ($1, $2)",
		trackingState.State, trackingState.UserID)
	if err != nil {
		log.Println("Error creating trackingState:", err)
		return err
	}
	return nil
}

func GetTrackingStatesByUserID(userID string) ([]TrackingState, error) {
	rows, err := mydb.GlobalDB.Query("SELECT id, amount, date, planned FROM trackingState WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error querying trackingStates:", err)
		return nil, err
	}
	defer rows.Close()

	var trackingStates []TrackingState
	for rows.Next() {
		var trackingState TrackingState
		if err := rows.Scan(&trackingState.State); err != nil {
			log.Println("Error scanning trackingState row:", err)
			return nil, err
		}
		trackingState.UserID = userID
		trackingStates = append(trackingStates, trackingState)
	}

	return trackingStates, nil
}
