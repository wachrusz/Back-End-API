package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/models"
	"github.com/wachrusz/Back-End-API/internal/service/categories"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/util"
	"go.uber.org/zap"
	"net/http"
)

type ProfileResponse struct {
	Message    string           `json:"message"`
	Profile    user.UserProfile `json:"profile"`
	StatusCode int              `json:"status_code"`
}

// GetProfileHandler retrieves the profile for the authenticated user.
//
// @Summary Get user profile
// @Description Get the user profile for the authenticated user.
// @Tags Profile
// @Produce json
// @Success 200 {object} ProfileResponse "User profile retrieved successfully"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error getting user profile"
// @Security JWT
// @Router /profile [get]
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
	response := ProfileResponse{
		Message:    "Successfully got the profile",
		Profile:    *userProfile,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("User profile retrieved successfully", zap.String("userID", userID))
}

type ProfileAnalyticsResponse struct {
	Message    string               `json:"message"`
	Analytics  categories.Analytics `json:"analytics"`
	Currency   string               `json:"currency"`
	StatusCode int                  `json:"status_code"`
}

// GetProfileAnalyticsHandler retrieves analytics data for the user's profile.
//
// @Summary Get profile analytics
// @Description This endpoint allows authenticated users to retrieve profile analytics data, filtered by date range and pagination parameters.
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param   X-Currency  header   string  true  "Currency code for analytics data (e.g., USD, EUR)"
// @Param   limit       query    int     false "Limit for pagination"
// @Param   offset      query    int     false "Offset for pagination"
// @Param   start_date  query    string  false "Start date for analytics data (YYYY-MM-DD)"
// @Param   end_date    query    string  false "End date for analytics data (YYYY-MM-DD)"
// @Success 200 {object} ProfileAnalyticsResponse "Successfully retrieved analytics data"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error while fetching analytics data"
// @Security JWT
// @Router /profile/analytics [get]
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
		h.errResp(w, fmt.Errorf("failed to get analytics data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := ProfileAnalyticsResponse{
		Message:    "Successfully got Analytics",
		Analytics:  *analytics,
		Currency:   currencyCode,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

type ProfileTrackerResponse struct {
	Message    string             `json:"message"`
	Tracker    categories.Tracker `json:"tracker"`
	Currency   string             `json:"currency"`
	StatusCode int                `json:"status_code"`
}

// GetProfileTrackerHandler retrieves tracker data for the user's profile.
//
// @Summary Get profile tracker
// @Description This endpoint allows authenticated users to retrieve tracker data, with optional pagination parameters.
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param   X-Currency  header   string  true  "Currency code for tracker data (e.g., USD, EUR)"
// @Param   limit       query    int     false "Limit for pagination"
// @Param   offset      query    int     false "Offset for pagination"
// @Success 200 {object} ProfileTrackerResponse "Successfully retrieved tracker data"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error while fetching tracker data"
// @Security JWT
// @Router /profile/tracker [get]
func (h *MyHandler) GetProfileTrackerHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	currencyCode := r.Header.Get("X-Currency")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	tracker, err := h.s.Categories.GetTrackerFromDB(userID, currencyCode, limitStr, offsetStr)
	if err != nil {
		h.errResp(w, errors.New("Failed to get tracker data: "+err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := ProfileTrackerResponse{
		Message:    "Successfully got tracker",
		Tracker:    *tracker,
		Currency:   currencyCode,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

type ProfileMoreResponse struct {
	Message    string          `json:"message"`
	More       categories.More `json:"more"`
	StatusCode int             `json:"status_code"`
}

// GetProfileMore retrieves additional profile data for the user.
//
// @Summary Get profile additional data
// @Description This endpoint allows authenticated users to retrieve additional profile data.
// @Tags Profile
// @Accept  json
// @Produce  json
// @Success 200 {object} ProfileMoreResponse "Successfully retrieved additional profile data"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error while fetching additional data"
// @Security JWT
// @Router /profile/more [get]
func (h *MyHandler) GetProfileMore(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	more, err := h.s.Categories.GetMoreFromDB(userID)
	if err != nil {
		h.errResp(w, fmt.Errorf("failed to get more data: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := ProfileMoreResponse{
		Message:    "Successfully got more",
		More:       *more,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

type ProfileArchiveResponse struct {
	Message    string             `json:"message"`
	Archive    []models.Operation `json:"archive"`
	StatusCode int                `json:"status_code"`
}

// GetOperationArchive retrieves the archived operations for the user's profile.
//
// @Summary Get operation archive
// @Description This endpoint allows authenticated users to retrieve their archived operations with optional pagination.
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param   limit   query    int     false "Limit for pagination"
// @Param   offset  query    int     false "Offset for pagination"
// @Success 200 {object} ProfileArchiveResponse "Successfully retrieved operation archive"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error while fetching operation archive"
// @Security JWT
// @Router /profile/archive [get]
func (h *MyHandler) GetOperationArchive(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	operations, err := h.s.Categories.GetOperationArchiveFromDB(userID, limitStr, offsetStr)
	if err != nil {
		h.errResp(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := ProfileArchiveResponse{
		Message:    "Successfully retrieved operation archive",
		Archive:    operations,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
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
// @Success 200 {object} jsonresponse.SuccessResponse "User profile updated successfully"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error updating user profile"
// @Security JWT
// @Router /profile/name [put]
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
	response := jsonresponse.SuccessResponse{
		Message:    "User profile updated successfully",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("User profile updated successfully", zap.String("userID", userID))
}
