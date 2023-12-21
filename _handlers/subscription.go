//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	jsonresponse "backEndAPI/_json_response"
	models "backEndAPI/_models"
	"encoding/json"
	"errors"
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

	err := models.CreateSubscription(&subscription)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating subscription: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully created a subscription",
		"status_code": http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}
