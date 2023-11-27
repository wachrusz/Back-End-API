//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	models "backEndAPI/_models"
	"encoding/json"
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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := models.CreateSubscription(&subscription)
	if err != nil {
		http.Error(w, "Error creating subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Subscription created successfully"))
}
