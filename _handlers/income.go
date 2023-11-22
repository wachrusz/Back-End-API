package handlers

import (
	auth "backEndAPI/_auth"
	models "backEndAPI/_models"
	"log"

	"encoding/json"
	"net/http"
)

func CreateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Started CreateIncomeHandler")
	var income models.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deviceID := auth.GetDeviceIDFromRequest(r)

	userID, ok := auth.GetUserIDFromSessionDatabase(deviceID)
	if ok != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	income.UserID = userID

	if err := models.CreateIncome(&income); err != nil {
		http.Error(w, "Error creating income", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Income created successfully"))
}
