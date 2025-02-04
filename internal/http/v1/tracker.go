package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
)

type GoalRequest struct {
	Goal models.Goal `json:"goal"`
}

// CreateGoalHandler creates a new goal in the database.
//
// @Summary Create a goal
// @Description Create a new goal.
// @Tags Tracker
// @Accept json
// @Produce json
// @Param goal body GoalRequest true "Goal object"
// @Success 201 {object} jsonresponse.IdResponse "Goal created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating goal"
// @Security JWT
// @Router /tracker/goal [post]
func (h *MyHandler) CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new goal...")

	// Decode the request payload
	var goalR GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&goalR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	goal := goalR.Goal

	// Extract the user ID from the request context
	userIDStr, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.errResp(w, fmt.Errorf("invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	// Assign user ID to the goal
	goal.UserID = userID

	// Create a new goal in the database
	goalID, err := h.m.Goals.Create(&goal)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating goal: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := jsonresponse.IdResponse{
		Message:    "Successfully created a goal",
		Id:         goalID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Goal created successfully", zap.Int64("goalID", goalID))
}

// UpdateGoalHandler updates an existing goal in the database.
//
// @Summary Update the goal
// @Description Updates an existing goal. There is no need to fill user_id field.
// @Tags Tracker
// @Accept json
// @Produce json
// @Param goal body GoalRequest true "Goal object"
// @Success 201 {object} jsonresponse.IdResponse "Goal updated successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 404 {object} jsonresponse.ErrorResponse "Goal not found"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error updating goal"
// @Security JWT
// @Router /tracker/goal [put]
func (h *MyHandler) UpdateGoalHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Updating a goal...")

	// Decode the request payload
	var goalR GoalRequest
	if err := json.NewDecoder(r.Body).Decode(&goalR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	goal := goalR.Goal

	// Extract the user ID from the request context
	userIDStr, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.errResp(w, fmt.Errorf("invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	// Assign user ID to the goal
	goal.UserID = userID

	// Create a new goal in the database
	if err := h.m.Goals.Update(&goal); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("expense not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error updating expense: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Send success response
	response := jsonresponse.IdResponse{
		Message:    "Successfully updated a goal",
		Id:         goal.ID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Goal updated successfully", zap.Int64("goalID", goal.ID))
}

// DeleteGoalHandler handles the deletion of an existing goal.
//
// @Summary Delete the goal
// @Description Delete the existing goal.
// @Tags Tracker
// @Param ConnectedAccount body jsonresponse.IdRequest true "goal id"
// @Success 204 {object} jsonresponse.SuccessResponse "goal deleted successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error deleting goal"
// @Security JWT
// @Router /tracker/goal [delete]
func (h *MyHandler) DeleteGoalHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting goal...")

	var id jsonresponse.IdRequest
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	goalID, err := strconv.ParseInt(id.ID, 10, 64)
	if err != nil {
		h.errResp(w, fmt.Errorf("invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	userIDStr, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.errResp(w, fmt.Errorf("invalid user ID: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.m.Goals.Delete(goalID, userID); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("goal not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error deleting goal: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted goal",
		StatusCode: http.StatusNoContent,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
