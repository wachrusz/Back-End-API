package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/util"
	"go.uber.org/zap"
	"net/http"
)

// GetProfileHandler retrieves the profile for the authenticated user.
//
// @Summary Get user profile
// @Description Get the user profile for the authenticated user.
// @Tags Profile
// @Produce json
// @Success 200 {string} string "User profile retrieved successfully"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error getting user profile"
// @Security JWT
// @Router /profile/get [get]
func (h *MyHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Retrieving user profile...")

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Fetch the user profile from the database
	userProfile, err := h.s.Users.GetProfile(userID)
	if err != nil {
		h.errResp(w, fmt.Errorf("error getting user profile: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response with user profile
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":     "Successfully got a profile",
		"status_code": http.StatusOK,
		"profile":     userProfile,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("User profile retrieved successfully", zap.String("userID", userID))
}

func (h *MyHandler) GetProfileAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	currencyCode := r.Header.Get("X-Currency")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	analytics, err := h.s.Categories.GetAnalyticsFromDB(userID, currencyCode, limitStr, offsetStr, startDateStr, endDateStr)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("failed to get analytics data: %v", err), http.StatusInternalServerError)
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

func (h *MyHandler) GetProfileTrackerHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	currencyCode := r.Header.Get("X-Currency")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	tracker, err_trk := h.s.Categories.GetTrackerFromDB(userID, currencyCode, limitStr, offsetStr)
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

func (h *MyHandler) GetProfileMore(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	more, err := h.s.Categories.GetMoreFromDB(userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("failed to get more data: %v", err), http.StatusInternalServerError)
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

func (h *MyHandler) GetOperationArchive(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	operations, err := h.s.Categories.GetOperationArchiveFromDB(userID, limitStr, offsetStr)
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

// UpdateName updates the authenticated user's profile with a new name and surname.
//
// @Summary Update user profile with name
// @Description Update the user profile for the authenticated user with a new name and surname.
// @Tags Profile
// @Accept json
// @Produce json
// @Param name body string true "New name to be added to the profile"
// @Success 200 {string} string "User profile updated successfully"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error updating user profile"
// @Security JWT
// @Router /profile/update-name [put]
func (h *MyHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Updating user profile...")

	// Extract user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Decode the request payload
	var request struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.errResp(w, fmt.Errorf("error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Update the user's name and surname in the database
	if err := h.s.Users.UpdateUserNameInDB(userID, request.Name, request.Surname); err != nil {
		h.errResp(w, fmt.Errorf("error updating name in the database: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"message":     "Successfully updated a profile",
		"status_code": http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("User profile updated successfully", zap.String("userID", userID))
}
