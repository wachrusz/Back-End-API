package v1

import (
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/models"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"go.uber.org/zap"
	"net/http"
)

// CreateWealthFundHandler creates a new wealth fund in the database.
//
// @Summary Create a wealth fund
// @Description Create a new wealth fund.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param wealthFund body models.WealthFund true "Wealth fund object"
// @Success 201 {string} string "Wealth fund created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error creating wealth fund"
// @Security JWT
// @Router /analytics/wealth_fund [post]
func (h *MyHandler) CreateWealthFundHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new wealth fund...")

	// Decode the request payload
	var wealthFund models.WealthFund
	if err := json.NewDecoder(r.Body).Decode(&wealthFund); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the wealth fund
	wealthFund.UserID = userID

	// Create a new wealth fund in the database
	wealthFundID, err := models.CreateWealthFund(&wealthFund)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating wealth fund: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"message":           "Successfully created a wealth fund",
		"created_object_id": wealthFundID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Wealth fund created successfully", zap.Int64("wealthFundID", wealthFundID))
}
