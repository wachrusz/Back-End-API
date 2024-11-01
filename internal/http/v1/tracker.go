package v1

import (
	"encoding/json"
	"fmt"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"go.uber.org/zap"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/models"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
)

// CreateGoalHandler creates a new goal in the database.
//
// @Summary Create a goal
// @Description Create a new goal.
// @Tags Tracker
// @Accept json
// @Produce json
// @Param goal body models.Goal true "Goal object"
// @Success 201 {object} jsonresponse.IdResponse "Goal created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating goal"
// @Security JWT
// @Router /tracker/goal [post]
func (h *MyHandler) CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new goal...")

	// Decode the request payload
	var goal models.Goal
	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the goal
	goal.UserID = userID

	// Create a new goal in the database
	goalID, err := models.CreateGoal(&goal)
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