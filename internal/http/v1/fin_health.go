package v1

/*
func CreateFinHealthHandler(w http.ResponseWriter, r *http.Request) {
	var finHealth models.FinHealth
	if errResp := json.NewDecoder(r.Body).Decode(&finHealth); errResp != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deviceID := auth.GetDeviceIDFromRequest(r)

	userID, ok := auth.GetUserIDFromSessionDatabase(deviceID)
	if ok != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	finHealth.UserID = userID

	if errResp := models.CreateFinHealth(&finHealth); errResp != nil {
		http.Error(w, "Error creating finHealth", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("FinHealth created successfully"))
}
*/
