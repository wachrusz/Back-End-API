//go:build !exclude_swagger
// +build !exclude_swagger

// Package profile provides profile information and it's functionality.
package profile

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/wachrusz/Back-End-API/internal/auth"
	"github.com/wachrusz/Back-End-API/internal/categories"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	mydb "github.com/wachrusz/Back-End-API/pkg/mydatabase"
	"net/http"
)

type UserProfile struct {
	Surname   string `json:"surname"` //*changed
	Name      string `json:"name"`
	UserID    string `json:"user_id"`
	AvatarURL string `json:"avatar_url"`
}

var (
	limitStr  string = "20"
	offsetStr string = "0"
)

func RegisterHandlers(router chi.Router) {
	// Profile routes
	router.Route("/profile", func(r chi.Router) {
		r.Get("/info/get", auth.AuthMiddleware(GetProfile))
		r.Get("/analytics/get", auth.AuthMiddleware(GetProfileAnalytics))
		r.Get("/tracker/get", auth.AuthMiddleware(GetProfileTracker))
		r.Get("/more/get", auth.AuthMiddleware(GetProfileMore))
		r.Put("/name/put", auth.AuthMiddleware(UpdateName))
		r.Get("/operation-archive/get", auth.AuthMiddleware(GetOperationArchive))
		r.Put("/image/put", auth.AuthMiddleware(UploadAvatarHandler))
	})

	// Emojis routes
	router.Route("/api/emojis", func(r chi.Router) {
		r.Put("/put", UploadIconHandler)
		r.Get("/get/list", GetIconsURLs)
	})
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
	var userProfile UserProfile
	surname, name, err := categories.GetUserInfoFromDB(userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Failed to get tracker data: "+err.Error()), http.StatusInternalServerError)
		return
	}
	avatarURL, err := GetAvatarInfo(userID)
	if err != nil {
		userProfile.AvatarURL = "null"
	}

	userProfile.UserID = userID
	userProfile.Surname = surname
	userProfile.Name = name
	userProfile.AvatarURL = avatarURL

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":     "Successfully got a profile",
		"status_code": http.StatusOK,
		"profile":     userProfile,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func GetProfileAnalytics(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	currencyCode := r.Header.Get("X-Currency")
	limitStr = r.URL.Query().Get("limit")
	offsetStr = r.URL.Query().Get("offset")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	analytics, err := categories.GetAnalyticsFromDB(userID, currencyCode, limitStr, offsetStr, startDateStr, endDateStr)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Failed to get analytics data: "+err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":           "Successfully got analytics",
		"status_code":       http.StatusOK,
		"analytics":         analytics,
		"response_currency": currencyCode,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func GetProfileTracker(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	currencyCode := r.Header.Get("X-Currency")
	limitStr = r.URL.Query().Get("limit")
	offsetStr = r.URL.Query().Get("offset")

	tracker, err_trk := categories.GetTrackerFromDB(userID, currencyCode, limitStr, offsetStr)
	if err_trk != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Failed to get tracker data: "+err_trk.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":           "Successfully got tracker",
		"status_code":       http.StatusOK,
		"tracker":           tracker,
		"response_currency": currencyCode,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func GetProfileMore(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	more, err := categories.GetMoreFromDB(userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Failed to get more data: "+err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":     "Successfully got more",
		"status_code": http.StatusOK,
		"more":        more,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func GetOperationArchive(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated"), http.StatusUnauthorized)
		return
	}

	limitStr = r.URL.Query().Get("limit")
	offsetStr = r.URL.Query().Get("offset")

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
	w.WriteHeader(response["status_code"].(int))
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
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func UpdateUserNameInDB(userID, newName, newSurname string) error {
	_, err := mydb.GlobalDB.Exec("UPDATE users SET name = $1, surname = $3 WHERE id = $2", newName, userID, newSurname)
	return err
}
