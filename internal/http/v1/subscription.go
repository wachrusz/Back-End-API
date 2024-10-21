package v1

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"

	"github.com/wachrusz/Back-End-API/internal/models"
	"net/http"
)

// CreateSubscriptionHandler creates a new subscription in the database.
//
// @Summary Create a subscription
// @Description Create a new subscription record.
// @Tags Settings
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription object"
// @Success 201 {string} string "Subscription created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating subscription"
// @Router /settings/subscription [post]
func (h *MyHandler) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new subscription...")

	// Decode the request payload
	var subscription models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Create a new subscription in the database
	subscriptionID, err := models.CreateSubscription(&subscription)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating subscription: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"message":           "Successfully created a subscription",
		"created_object_id": subscriptionID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Subscription created successfully", zap.Int64("subscriptionID", subscriptionID))
}
