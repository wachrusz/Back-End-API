package handlers

import (
	auth "backEndAPI/_auth"
	models "backEndAPI/_models"

	"encoding/json"
	"net/http"
)

func CreateWealthFundHandler(w http.ResponseWriter, r *http.Request) {
	var wealthFund models.WealthFund
	if err := json.NewDecoder(r.Body).Decode(&wealthFund); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deviceID := auth.GetDeviceIDFromRequest(r)

	userID, ok := auth.GetUserIDFromSessionDatabase(deviceID)
	if ok != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	wealthFund.UserID = userID

	if err := models.CreateWealthFund(&wealthFund); err != nil {
		http.Error(w, "Error creating wealthFund", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("WealthFund created successfully"))
}
