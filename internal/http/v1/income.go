package v1

import (
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/models"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"go.uber.org/zap"
	"net/http"
)

// CreateIncomeHandler creates a new income record in the database.
//
// @Summary Create an income
// @Description Create a new income record.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param income body models.Income true "Income object"
// @Success 201 {string} string "Income created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error creating income"
// @Security JWT
// @Router /analytics/income [post]
func (h *MyHandler) CreateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new income...")

	// Decode the request payload
	var income models.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the income
	income.UserID = userID

	// Create a new income in the database
	incomeID, err := models.CreateIncome(&income)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating income: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"message":           "Successfully created an income",
		"created_object_id": incomeID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Income created successfully", zap.Int64("incomeID", incomeID))
}
