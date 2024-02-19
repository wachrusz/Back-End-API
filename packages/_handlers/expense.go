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
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	expense.UserID = userID

	if err := models.CreateExpense(&expense); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating expense: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully created an expense",
		"status_code": http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
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
