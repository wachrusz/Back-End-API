//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package v1

import (
	"encoding/json"
	"errors"
	"github.com/wachrusz/Back-End-API/internal/models"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
)

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
	var wealthFund models.WealthFund
	if err := json.NewDecoder(r.Body).Decode(&wealthFund); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	wealthFund.UserID = userID

	wealthFundID, err := models.CreateWealthFund(&wealthFund)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating wealthFund: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":           "Successfully created a wealth fund",
		"created_object_id": wealthFundID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}
