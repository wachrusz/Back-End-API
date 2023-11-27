//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "backEndAPI/_mydatabase"
	"log"
)

type TrackingState struct {
	State  float64 `json:"state"`
	UserID string  `json:"user_id"`
}

// @Summary Create tracking state
// @Description Create a new tracking state entry.
// @Tags TrackingState
// @Accept json
// @Produce json
// @Param trackingState body TrackingState true "Tracking state details"
// @Success 201 {string} string "Tracking state created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating tracking state"
// @Router /models/tracking-state [post]
func CreateTrackingState(trackingState *TrackingState) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO trackingState (state, user_id) VALUES ($1, $2)",
		trackingState.State, trackingState.UserID)
	if err != nil {
		log.Println("Error creating trackingState:", err)
		return err
	}
	return nil
}

// @Summary Get tracking states by user ID
// @Description Get a list of tracking states for a specific user.
// @Tags TrackingState
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {array} TrackingState "List of tracking states"
// @Failure 500 {string} string "Error querying tracking states"
// @Router /models/tracking-state/{userID} [get]
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
