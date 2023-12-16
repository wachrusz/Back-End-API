//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	auth "backEndAPI/_auth"
	models "backEndAPI/_models"

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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	wealthFund.UserID = userID

	if err := models.CreateWealthFund(&wealthFund); err != nil {
		http.Error(w, "Error creating wealthFund", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("WealthFund created successfully"))
}
