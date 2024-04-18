//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	jsonresponse "main/packages/_json_response"
	models "main/packages/_models"
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
// @Security JWT
// @Router /app/category/expense [post]
func CreateExpenseCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.ExpenseCategory

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	err := models.CreateExpenseCategory(&category)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating expense category: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Expense category created successfully",
		"status_code": http.StatusCreated,
	}
	w.WriteHeader(response["status_code"])
	json.NewEncoder(w).Encode(response)
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
// @Security JWT
// @Router /app/category/income [post]
func CreateIncomeCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.IncomeCategory

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	err := models.CreateIncomeCategory(&category)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating income category: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Income category created successfully",
		"status_code": http.StatusCreated,
	}
	w.WriteHeader(response["status_code"])
	json.NewEncoder(w).Encode(response)
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
// @Security JWT
// @Router /app/category/investment [post]
func CreateInvestmentCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.InvestmentCategory

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	err := models.CreateInvestmentCategory(&category)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating investment category: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Investment category created successfully",
		"status_code": http.StatusCreated,
	}
	w.WriteHeader(response["status_code"])
	json.NewEncoder(w).Encode(response)
}
