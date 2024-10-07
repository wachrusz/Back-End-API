package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/util"
	"github.com/wachrusz/Back-End-API/pkg/validator"
	"log"
	"net/http"
	"time"
)

// @Summary Register user
// @Description Register a new user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Param name query string true "Name"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error registering user"
// @Router /auth/register [post]
func (h *MyHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var registrationRequest user.UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&registrationRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid request payload: %v", err.Error()), http.StatusBadRequest)
		return
	}

	// validation

	if !validator.IsValidEmail(registrationRequest.Email) {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid email: %s", registrationRequest.Email), http.StatusBadRequest)
		return
	}

	if !validator.IsValidPassword(registrationRequest.Password) {
		log.Println("invalid password")
		jsonresponse.SendErrorResponse(w, fmt.Errorf("password must be at least 7 digits long"), http.StatusBadRequest)
		return
	}

	token, err := h.s.Users.PrimaryRegistration(registrationRequest.Email, registrationRequest.Password)
	if err != nil {
		switch err {
		case myerrors.ErrDuplicated:

		default:
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid request payload: %v", err.Error()), http.StatusBadRequest)
		}
		return
	}

	response := map[string]interface{}{
		"message":     "Confirm your email",
		"token":       token,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MyHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var resetRequest user.ResetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	if !validator.IsValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	err = h.s.Users.ResetPassword(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
	}
}

func (h *MyHandler) ChangePasswordForRecoverHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var resetRequest user.UserPasswordReset

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	password := resetRequest.Password
	if !validator.IsValidPassword(password) {
		jsonresponse.SendErrorResponse(w, errors.New("password must be at least 7 digits long: "), http.StatusBadRequest)
		return
	}

	if !validator.IsValidEmail(email) {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	resetToken := resetRequest.ResetToken
	err = h.s.Users.ChangePasswordForRecover(email, password, resetToken)
	if err != nil {
		var statusCode = 500
		switch {
		case errors.Is(err, myerrors.ErrEmpty):
			statusCode = http.StatusBadRequest
			break
		case errors.Is(err, myerrors.ErrEmailing):
			statusCode = http.StatusInternalServerError
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
			statusCode = http.StatusUnauthorized
			break
		}
		jsonresponse.SendErrorResponse(w, fmt.Errorf("error changing the password: %v", err), statusCode)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfuly reseted password",
		"status_code": http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Login to the system
// @Description Login to the system and get an authentication token
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body auth.UserAuthenticationRequest true "UserAuthenticationRequest object"
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/login [post]
func (h *MyHandler) Login(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}

	var loginRequest user.UserAuthenticationRequest

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := loginRequest.Email
	password := loginRequest.Password

	token, err := h.s.Users.Login(email, password)
	if err != nil {
		switch {
		case errors.Is(err, myerrors.ErrEmpty) || errors.Is(err, myerrors.ErrInvalidCreds):
			jsonresponse.SendErrorResponse(w, fmt.Errorf("invalid email or password: %w", err), http.StatusUnauthorized)
		case errors.Is(err, myerrors.ErrInternal) || errors.Is(err, myerrors.ErrEmailing):
			jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		}
		return
	}

	response := map[string]interface{}{
		"message":     "Confirm your email",
		"token":       token,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func (h *MyHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}
	//! Заставляет задуматься
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	//! Сомнительно
	type token struct {
		RefreshToken string `json:"refresh_token"`
	}

	var tmpToken token

	err := json.NewDecoder(r.Body).Decode(&tmpToken)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	access, refresh, err := h.s.Users.RefreshToken(tmpToken.RefreshToken, userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":                 "Successfully refreshed tokens",
		"access_token":            access,
		"refresh_token":           refresh,
		"access_token_life_time":  time.Minute * 15,
		"refresh_token_life_time": 30 * 24 * time.Hour,
		"status_code":             http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

// @Summary Logout the user
// @Description Logs out the user, terminating the session.
// @Tags Auth
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Failure 500 {string} string "Internal Server Error"
// @Security JWT
// @Router /auth/logout [post]
func (h *MyHandler) Logout(w http.ResponseWriter, r *http.Request) {
	currentDeviceID, ok := utility.GetDeviceIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	err := h.s.Users.Logout(currentDeviceID, userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Logout Successful",
		"status_code": http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
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
func (h *MyHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	userProfile, err := h.s.Users.GetProfile(userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":     "Successfully got a profile",
		"status_code": http.StatusOK,
		"profile":     userProfile,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
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
func (h *MyHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
	userID, ok := utility.GetUserIDFromContext(r.Context())
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

	err := h.s.Users.UpdateUserNameInDB(userID, request.Name, request.Surname)
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

func (h *MyHandler) DeleteTokensHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	deviceID := r.URL.Query().Get("deviceID")
	if (email == "" && deviceID == "") || (email != "" && deviceID != "") {
		jsonresponse.SendErrorResponse(w, errors.New("blank fields and two methods are not allowed"), http.StatusBadRequest)
		return
	}

	err := h.s.Users.DeleteTokens(email, deviceID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":     "Successfuly deleted tokens",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *MyHandler) GetTokenPairsAmountHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("blank fields are not allowed"), http.StatusBadRequest)
		return
	}
	amount, err := h.s.Users.GetTokenPairsAmount(email)
	if err != nil {
		jsonresponse.SendErrorResponse(w, fmt.Errorf("error while counting sessions: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":     "Successfuly got ammount",
		"ammount":     amount,
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}
