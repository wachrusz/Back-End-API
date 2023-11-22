package handlers

import (
	auth "backEndAPI/_auth"
	models "backEndAPI/_models"
	"log"

	"encoding/json"
	"net/http"
)

func CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Started CreateGoalHandler")
	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deviceID := auth.GetDeviceIDFromRequest(r)

	userID, ok := auth.GetUserIDFromSessionDatabase(deviceID)
	if ok != nil {
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
