package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"go.uber.org/zap"
	"net/http"
)

type SubscriptionRequest struct {
	Subscription models.Subscription `json:"subscription"`
}

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
// @Param subscription body SubscriptionRequest true "Subscription object"
// @Success 201 {object} EndTimeResponse "Subscription created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating subscription"
// @Security JWT
// @Router /settings/subscription [post]
func (h *MyHandler) CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new subscription...")

	// Decode the request payload
	var subscriptionR SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subscriptionR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	subscription := subscriptionR.Subscription

	// Create a new subscription in the database
	subscriptionID, err := h.m.Subscriptions.Create(&subscription)
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

// UpdateSubscriptionHandler updates an existing subscription in the database.
//
// @Summary Update the subscription
// @Description Updates an existing subscription. There is no need to fill user_id field.
// @Tags Settings
// @Accept json
// @Produce json
// @Param subscription body SubscriptionRequest true "Subscription object"
// @Success 201 {object} jsonresponse.IdResponse "subscription updated successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 404 {object} jsonresponse.ErrorResponse "Subscription not found"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error updating subscription"
// @Security JWT
// @Router /settings/subscription [put]
func (h *MyHandler) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Updating a subscription...")

	// Decode the request payload
	var subscriptionR SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subscriptionR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	subscription := subscriptionR.Subscription

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the subscription
	subscription.UserID = userID

	// Create a new subscription in the database
	if err := h.m.Subscriptions.Update(&subscription); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("expense not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error updating expense: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Send success response
	response := jsonresponse.IdResponse{
		Message:    "Successfully updated a subscription",
		Id:         subscription.ID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Subscription updated successfully", zap.Int64("subscriptionID", subscription.ID))
}

// DeleteSubscriptionHandler handles the deletion of an existing subscription.
//
// @Summary Delete the subscription
// @Description Delete the existing subscription.
// @Tags Settings
// @Param ConnectedAccount body jsonresponse.IdRequest true "Subscription id"
// @Success 204 {object} jsonresponse.SuccessResponse "Subscription deleted successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error deleting subscription"
// @Security JWT
// @Router /settings/subscription [delete]
func (h *MyHandler) DeleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting subscription...")

	var id jsonresponse.IdRequest
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	if err := h.m.Subscriptions.Delete(id.ID, userID); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("subscription not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error deleting subscription: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted subscription",
		StatusCode: http.StatusNoContent,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
