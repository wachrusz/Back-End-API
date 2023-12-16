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

// @Summary Create an expense
// @Description Create a new expense.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param expense body models.Expense true "Expense object"
// @Success 201 {string} string "Expense created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error creating expense"
// @Security JWT
// @Router /analytics/expense [post]
func CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	expense.UserID = userID

	if err := models.CreateExpense(&expense); err != nil {
		http.Error(w, "Error creating expense", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Expense created successfully"))
}

/*
func GetExpensesHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.Login(r)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expenses, err := models.GetExpensesByUserID(userID)
	if err != nil {
		http.Error(w, "Error getting expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}
*/
