//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	"errors"
	"log"

	auth "main/packages/_auth"
	jsonresponse "main/packages/_json_response"
	models "main/packages/_models"

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
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	goal.UserID = userID
	log.Print(goal.UserID)

	goalID, err := models.CreateGoal(&goal)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating goal: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":           "Successfully created a goal",
		"created_object_id": goalID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}
