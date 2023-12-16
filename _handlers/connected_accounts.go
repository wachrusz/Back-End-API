//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	auth "backEndAPI/_auth"
	models "backEndAPI/_models"
	"encoding/json"
	"net/http"
)

// @Summary Create a connected account
// @Description Create a new connected account.
// @Tags App
// @Accept json
// @Produce json
// @Param ConnectedAccount body models.ConnectedAccount true "ConnectedAccount object"
// @Success 201 {string} string "Connected account created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error adding connected account"
// @Security JWT
// @Router /app/connected-accounts/add [post]
func AddConnectedAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account models.ConnectedAccount

	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := models.AddConnectedAccount(&account)
	if err != nil {
		http.Error(w, "Error adding connected account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Connected account added successfully"))
}

// @Summary Delete a connected account
// @Description Delete an existing connected account.
// @Tags App
// @Param ConnectedAccount body models.ConnectedAccount true "ConnectedAccount object"
// @Success 201 {string} string "Connected account created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error adding connected account"
// @Security JWT
// @Router /app/connected-accounts/delete [delete]
func DeleteConnectedAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Error deleting connected account: UNAUTHORIZED", http.StatusUnauthorized)
		return
	}

	err := models.DeleteConnectedAccount(userID)
	if err != nil {
		http.Error(w, "Error deleting connected account: DB_Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Connected account deleted successfully"))
}
