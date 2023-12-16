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

// @Summary Create an income
// @Description Create a new income.
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
func CreateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	var income models.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	income.UserID = userID

	if err := models.CreateIncome(&income); err != nil {
		http.Error(w, "Error creating income", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Income created successfully"))
}
