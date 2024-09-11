//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package v1

import (
	"encoding/json"
	"errors"
	"github.com/wachrusz/Back-End-API/internal/auth"
	"github.com/wachrusz/Back-End-API/internal/models"
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
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	account.UserID = userID

	connectedAccountID, err := models.AddConnectedAccount(&account)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error adding connected account: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":           "Connected account added successfully",
		"created_object_id": connectedAccountID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
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
		jsonresponse.SendErrorResponse(w, errors.New("Error deleting connected account: UNAUTHORIZED"), http.StatusUnauthorized)
		return
	}

	err := models.DeleteConnectedAccount(userID)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error deleting connected account: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfuly deleted connected account",
		"status_code": http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}
