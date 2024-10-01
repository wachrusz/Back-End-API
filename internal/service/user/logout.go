//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package user

import (
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/auth/service"
)

// @Summary Logout the user
// @Description Logs out the user, terminating the session.
// @Tags Auth
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Failure 500 {string} string "Internal Server Error"
// @Security JWT
// @Router /auth/logout [post]
func Logout(device, userID string) error {
	err := service.RemoveSessionFromDatabase(device, userID)
	if err != nil {
		return fmt.Errorf("error removing session from db: %v", err)
	}

	delete(service.ActiveUsers, userID)
	return nil
}
