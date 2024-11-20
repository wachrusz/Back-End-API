package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"go.uber.org/zap"
	"net/http"

	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
)

type ExpenseRequest struct {
	Expense models.Expense `json:"expense"`
}

// CreateExpenseHandler creates a new expense record in the database.
//
// @Summary Create a expense
// @Description Create a new expense record.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param expense body ExpenseRequest true "Expense object"
// @Success 201 {object} jsonresponse.IdResponse "Successfully created an expense"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating expense"
// @Security JWT
// @Router /analytics/expenses [post]
func (h *MyHandler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new expense...")

	// Decode the request payload
	var expenseR ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&expenseR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	expense := expenseR.Expense

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the expense
	expense.UserID = userID

	// Create a new expense in the database
	expenseID, err := h.m.Expenses.Create(&expense)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating expense: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := jsonresponse.IdResponse{
		Message:    "Successfully created an expense",
		Id:         expenseID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Expense created successfully", zap.Int64("expenseID", expenseID))
}

// UpdateExpenseHandler handles the update of an existing expense.
//
// @Summary Update the expense
// @Description Update an existing expense. There is no need to fill user_id field.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param ConnectedAccount body ExpenseRequest true "Expense object"
// @Success 200 {object} jsonresponse.SuccessResponse "expense updated successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 404 {object} jsonresponse.ErrorResponse "expense not found"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error updating expense"
// @Security JWT
// @Router /analytics/expense [put]
func (h *MyHandler) UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Updating expense...")

	// Decode the request body
	var expenseR ExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&expenseR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	expense := expenseR.Expense

	// Check if user is authenticated
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}
	expense.UserID = userID

	// Attempt to update the account
	if err := h.m.Expenses.Update(&expense); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("expense not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error updating expense: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	response := jsonresponse.SuccessResponse{
		Message:    "expense updated successfully",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// DeleteExpenseHandler handles the deletion of an existing expense.
//
// @Summary Delete the expense
// @Description Delete the existing expense.
// @Tags Analytics
// @Param ConnectedAccount body jsonresponse.IdRequest true "Expense id"
// @Success 204 {object} jsonresponse.SuccessResponse "Expense deleted successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error deleting expense"
// @Security JWT
// @Router /analytics/expense [delete]
func (h *MyHandler) DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting expense...")

	var id jsonresponse.IdRequest
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	if err := h.m.Expenses.Delete(id.ID, userID); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("expense not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error deleting expense: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted expense",
		StatusCode: http.StatusNoContent,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// CreateIncomeHandler creates a new income record in the database.
//
// @Summary Create an income
// @Description Create a new income record.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param income body repository.Income true "Income object"
// @Success 201 {object} jsonresponse.IdResponse "Income created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating income"
// @Security JWT
// @Router /analytics/income [post]
func (h *MyHandler) CreateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new income...")

	// Decode the request payload
	var income repository.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the income
	income.UserID = userID

	// Create a new income in the database
	incomeID, err := repository.CreateIncome(&income)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating income: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := jsonresponse.IdResponse{
		Message:    "Successfully created an income",
		Id:         incomeID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Income created successfully", zap.Int64("incomeID", incomeID))
}

// CreateWealthFundHandler creates a new wealth fund in the database.
//
// @Summary Create a wealth fund
// @Description Create a new wealth fund.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param wealthFund body repository.WealthFund true "Wealth fund object"
// @Success 201 {object} jsonresponse.IdResponse "Wealth fund created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating wealth fund"
// @Security JWT
// @Router /analytics/wealth_fund [post]
func (h *MyHandler) CreateWealthFundHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new wealth fund...")

	// Decode the request payload
	var wealthFund repository.WealthFund
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
	wealthFundID, err := repository.CreateWealthFund(&wealthFund)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating wealth fund: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := jsonresponse.IdResponse{
		Message:    "Successfully created a wealth fund",
		Id:         wealthFundID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Wealth fund created successfully", zap.Int64("wealthFundID", wealthFundID))
}
