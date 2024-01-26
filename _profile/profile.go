//go:build !exclude_swagger
// +build !exclude_swagger

// Package profile provides profile information and it's functionality.
package profile

import (
	"encoding/json"
	"errors"
	"net/http"

	auth "backEndAPI/_auth"
	categories "backEndAPI/_categories"
	jsonresponse "backEndAPI/_json_response"
	models "backEndAPI/_models"
	mydb "backEndAPI/_mydatabase"

	"github.com/gorilla/mux"
)

type UserProfile struct {
	Surname   string               `json:"surname"` //*changed
	Name      string               `json:"name"`
	Analytics categories.Analytics `json:"analytics"`
	Tracker   categories.Tracker   `json:"tracker"`
	More      categories.More      `json:"more"`
	UserID    string               `json:"userID"`
}

var userProfiles = make(map[string]UserProfile)

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/profile/get", auth.AuthMiddleware(GetProfile)).Methods("GET")
	router.HandleFunc("/profile/update-name", auth.AuthMiddleware(UpdateName)).Methods("PUT")
	router.HandleFunc("/profile/operation-archive/get", auth.AuthMiddleware(GetOperationArchive)).Methods("GET")

	router.HandleFunc("/profile/image/put", auth.AuthMiddleware(UploadAvatarHandler)).Methods("PUT")
	router.HandleFunc("/profile/image/get", auth.AuthMiddleware(GetAvatarHandler)).Methods("GET")
}

// @Summary Get user profile
// @Description Get the user profile for the authenticated user.
// @Tags Profile
// @Produce json
// @Success 200 {string} string "User profile retrieved successfully"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error getting user profile"
// @Security JWT
// @Router /profile/get [get]
func GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	currencyCode := r.Header.Get("X-Currency")

	var userProfile UserProfile

	userProfile.Analytics = categories.Analytics{
		Income:     make([]models.Income, 0),
		Expense:    make([]models.Expense, 0),
		WealthFund: make([]models.WealthFund, 0),
	}

	analytics, err := categories.GetAnalyticsFromDB(userID, currencyCode)
	if err != nil {
		http.Error(w, "Failed to get analytics data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tracker, err_trk := categories.GetTrackerFromDB(userID, currencyCode, analytics)
	if err_trk != nil {
		http.Error(w, "Failed to get tracker data: "+err_trk.Error(), http.StatusInternalServerError)
		return
	}
	surname, name, err := categories.GetUserInfoFromDB(userID)
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
	userProfile.Surname = surname
	userProfile.Name = name
	userProfile.Analytics = *analytics
	userProfile.Tracker = *tracker
	userProfile.More = *more

	userProfiles[userID] = userProfile

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":     "Successfully got a profile",
		"status_code": http.StatusOK,
		"profile":     userProfile,
	}
	json.NewEncoder(w).Encode(response)
}

func GetOperationArchive(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated"), http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr == "" || offsetStr == "" {
		limitStr = "20"
		offsetStr = "0"
	}

	operations, err := categories.GetOperationArchiveFromDB(userID, limitStr, offsetStr)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":           "Successfully got an archive",
		"status_code":       http.StatusOK,
		"operation_archive": operations,
	}
	json.NewEncoder(w).Encode(response)
}

// * Добавлены поля для имени и фамилии
// @Summary Update user profile with name
// @Description Update the user profile for the authenticated user with a new name.
// @Tags Profile
// @Accept json
// @Produce json
// @Param name body string true "New name to be added to the profile"
// @Success 200 {string} string "User profile updated successfully"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error updating user profile"
// @Security JWT
// @Router /profile/update-name [put]
func UpdateName(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	var request struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error decoding JSON: "+err.Error()), http.StatusBadRequest)
		return
	}

	err := UpdateUserNameInDB(userID, request.Name, request.Surname)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error updating name in the database: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully updated a profile",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUserNameInDB(userID, newName, newSurname string) error {
	_, err := mydb.GlobalDB.Exec("UPDATE users SET name = $1, surname = $3 WHERE id = $2", newName, userID, newSurname)
	return err
}
