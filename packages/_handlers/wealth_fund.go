//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	"errors"

	auth "main/packages/_auth"
	jsonresponse "main/packages/_json_response"
	models "main/packages/_models"

	"encoding/json"
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
func CreateWealthFundHandler(w http.ResponseWriter, r *http.Request) {
	var wealthFund models.WealthFund
	if err := json.NewDecoder(r.Body).Decode(&wealthFund); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	wealthFund.UserID = userID

	if err := models.CreateWealthFund(&wealthFund); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating wealthFund: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully created a wealth fund",
		"status_code": http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}
