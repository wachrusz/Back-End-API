//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	auth "backEndAPI/_auth"
	models "backEndAPI/_models"
	"log"

	"encoding/json"
	"net/http"
)

// @Summary Create a goal
// @Description Create a new goal.
// @Tags Tracker
// @Accept json
// @Produce json
// @Param goal body models.Goal true "Goal object"
// @Success 201 {string} string "Goal created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error creating goal"
// @Router /tracker/goal [post]
func CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Started CreateGoalHandler")
	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	goal.UserID = userID
	log.Print(goal.UserID)

	if err := models.CreateGoal(&goal); err != nil {
		http.Error(w, "Error creating goal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Goal created successfully"))
}
