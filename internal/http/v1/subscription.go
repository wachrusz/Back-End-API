//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	"encoding/json"
	"errors"
	"main/internal/models"
	"main/pkg/json_response"
	"net/http"
)

// @Summary Create a subscription
// @Description Create a new subscription.
// @Tags Settings
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription object"
// @Success 201 {string} string "Subscription created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating subscription"
// @Router /settings/subscription [post]
func CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var subscription models.Subscription

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	subscriptionID, err := models.CreateSubscription(&subscription)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating subscription: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":           "Successfully created a subscription",
		"created_object_id": subscriptionID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}
