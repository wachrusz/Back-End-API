//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	service "main/packages/_auth/service"
	jsonresponse "main/packages/_json_response"
)

// @Summary Logout the user
// @Description Logs out the user, terminating the session.
// @Tags Auth
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Failure 500 {string} string "Internal Server Error"
// @Security JWT
// @Router /auth/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {
	currentDeviceID, ok := GetDeviceIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	err := service.RemoveSessionFromDatabase(currentDeviceID, userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error removing session from the database: "+err.Error()), http.StatusInternalServerError)
		return
	}

	delete(service.ActiveUsers, userID)

	response := map[string]interface{}{
		"message":     "Logout Successful",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}
