//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	logger "backEndAPI/_logger"
	"net/http"
)

// @Summary Logout the user
// @Description Logs out the user, terminating the session.
// @Tags Auth
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Failure 500 {string} string "Internal Server Error"
// @Router /auth/logout [post]
func Logout(w http.ResponseWriter, r *http.Request) {

	currentDeviceID := GetDeviceIDFromRequest(r)

	userID, err := GetUserIDFromSessionDatabase(currentDeviceID)
	if err != nil {
		logger.ErrorLogger.Printf("Error getting userID from the database: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = removeSessionFromDatabase(currentDeviceID)
	if err != nil {
		logger.ErrorLogger.Printf("Error removing session from the database: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	delete(activeUsers, userID)

	logger.InfoLogger.Printf("User %s logged out from %s\n", userID, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
