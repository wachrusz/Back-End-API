package profile

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "backEndAPI/_auth"
	categories "backEndAPI/_categories"
	models "backEndAPI/_models"

	"github.com/gorilla/mux"
)

type UserProfile struct {
	Username  string               `json:"username"`
	Name      string               `json:"name"`
	Analytics categories.Analytics `json:"analytics"`
	Tracker   categories.Tracker   `json:"tracker"`
	More      categories.More      `json:"more"`
	UserID    string               `json:"userID"`
}

var userProfiles = make(map[string]UserProfile)

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/profile/get", GetProfile).Methods("GET")
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	deviceID := auth.GetDeviceIDFromRequest(r)

	userID, ok := auth.GetUserIDFromSessionDatabase(deviceID)
	if ok != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var userProfile UserProfile

	userProfile.Analytics = categories.Analytics{
		Income:     make([]models.Income, 0),
		Expense:    make([]models.Expense, 0),
		WealthFund: make([]models.WealthFund, 0),
	}

	analytics, err := categories.GetAnalyticsFromDB(userID)
	if err != nil {
		http.Error(w, "Failed to get analytics data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tracker, err_trk := categories.GetTrackerFromDB(userID, analytics)
	if err_trk != nil {
		http.Error(w, "Failed to get tracker data: "+err_trk.Error(), http.StatusInternalServerError)
		return
	}
	userName, name, err := categories.GetUserInfoFromDB(userID)
	if err != nil {
		http.Error(w, "Failed to get user data: "+err_trk.Error(), http.StatusInternalServerError)
		return
	}
	more, err := categories.GetMoreFromDB(userID)
	if err != nil {
		http.Error(w, "Failed to get More data: "+err_trk.Error(), http.StatusInternalServerError)
		return
	}

	userProfile.UserID = userID
	userProfile.Username = userName
	userProfile.Name = name
	userProfile.Analytics = *analytics
	userProfile.Tracker = *tracker
	userProfile.More = *more

	userProfiles[userID] = userProfile

	fmt.Printf("User Profile: %+v\n", userProfile)

	userProfileJSON, err := json.Marshal(userProfile)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userProfileJSON)
}
