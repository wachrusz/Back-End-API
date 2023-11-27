//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	models "backEndAPI/_models"
	"encoding/json"
	"net/http"
)

// @Summary Create an expense category
// @Description Create a new expense category.
// @Tags App
// @Accept json
// @Produce json
// @Param category body models.ExpenseCategory true "Expense category object"
// @Success 201 {string} string "Expense category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating expense category"
// @Router /app/category/expense [post]
func CreateExpenseCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.ExpenseCategory

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := models.CreateExpenseCategory(&category)
	if err != nil {
		http.Error(w, "Error creating expense category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Expense category created successfully"))
}

// @Summary Create an income category
// @Description Create a new income category.
// @Tags App
// @Accept json
// @Produce json
// @Param category body models.IncomeCategory true "Income category object"
// @Success 201 {string} string "Income category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating income category"
// @Router /app/category/income [post]
func CreateIncomeCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.IncomeCategory

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := models.CreateIncomeCategory(&category)
	if err != nil {
		http.Error(w, "Error creating income category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Income category created successfully"))
}

// @Summary Create an investment category
// @Description Create a new investment category.
// @Tags App
// @Accept json
// @Produce json
// @Param category body models.InvestmentCategory true "Investment category object"
// @Success 201 {string} string "Investment category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating investment category"
// @Router /app/category/investment [post]
func CreateInvestmentCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.InvestmentCategory

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := models.CreateInvestmentCategory(&category)
	if err != nil {
		http.Error(w, "Error creating income category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Investment category created successfully"))
}
