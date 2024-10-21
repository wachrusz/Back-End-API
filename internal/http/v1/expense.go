package v1

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/models"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
)

// CreateExpenseHandler creates a new expense record in the database.
//
// @Summary Create an expense
// @Description Create a new expense record.
// @Tags Expenses
// @Accept json
// @Produce json
// @Param expense body models.Expense true "Expense object"
// @Success 201 {string} string "Successfully created an expense"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error creating expense"
// @Router /expenses [post]
func (h *MyHandler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new expense...")

	// Decode the request payload
	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the expense
	expense.UserID = userID

	// Create a new expense in the database
	expenseID, err := models.CreateExpense(&expense)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating expense: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"message":           "Successfully created an expense",
		"created_object_id": expenseID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Expense created successfully", zap.Int64("expenseID", expenseID))
}
