package v1

import (
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/models"
	"go.uber.org/zap"
	"net/http"
)

type EndTimeResponse struct {
	Message    string `json:"message"`
	Id         int64  `json:"id"`
	EndTime    string `json:"end_date"`
	StatusCode int    `json:"status_code"`
}

// CreateSubscriptionHandler creates a new subscription in the database.
//
// @Summary Create a subscription
// @Description Create a new subscription record.
// @Tags Settings
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription object"
// @Success 201 {object} EndTimeResponse "Subscription created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating subscription"
// @Security JWT
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
	response := EndTimeResponse{
		Message:    "Successfully created a subscription",
		Id:         subscriptionID,
		EndTime:    subscription.EndDate,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Subscription created successfully", zap.Int64("subscriptionID", subscriptionID))
}
