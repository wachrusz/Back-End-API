package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"go.uber.org/zap"
	"net/http"

	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
)

// ExpenseRequest is used for deserialization
type ExpenseRequest struct {
	Expense models.Expense `json:"expense"`
}

// IncomeRequest is used for deserialization
type IncomeRequest struct {
	Income models.Income `json:"income"`
}

// WealthFundRequest is used for deserialization
type WealthFundRequest struct {
	WealthFund models.WealthFund `json:"wealth_fund"`
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
// @Param income body IncomeRequest true "Income object"
// @Success 201 {object} jsonresponse.IdResponse "Income created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating income"
// @Security JWT
// @Router /analytics/income [post]
func (h *MyHandler) CreateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new income...")

	// Decode the request payload
	var incomeR IncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&incomeR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	income := incomeR.Income

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the income
	income.UserID = userID

	// Create a new income in the database
	incomeID, err := h.m.Incomes.Create(&income)
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

// UpdateIncomeHandler handles the update of an existing expense.
//
// @Summary Update the income
// @Description Update an existing income. There is no need to fill user_id field.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param ConnectedAccount body IncomeRequest true "income object"
// @Success 200 {object} jsonresponse.SuccessResponse "income updated successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 404 {object} jsonresponse.ErrorResponse "income not found"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error updating income"
// @Security JWT
// @Router /analytics/income [put]
func (h *MyHandler) UpdateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Updating income...")

	// Decode the request body
	var incomeR IncomeRequest
	if err := json.NewDecoder(r.Body).Decode(&incomeR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	income := incomeR.Income

	// Check if user is authenticated
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}
	income.UserID = userID

	// Attempt to update the account
	if err := h.m.Incomes.Update(&income); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("income not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error updating income: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	response := jsonresponse.SuccessResponse{
		Message:    "income updated successfully",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// DeleteIncomeHandler handles the deletion of an existing income.
//
// @Summary Delete the income
// @Description Delete the existing income.
// @Tags Analytics
// @Param ConnectedAccount body jsonresponse.IdRequest true "income id"
// @Success 204 {object} jsonresponse.SuccessResponse "income deleted successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error deleting income"
// @Security JWT
// @Router /analytics/income [delete]
func (h *MyHandler) DeleteIncomeHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting income...")

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

	if err := h.m.Incomes.Delete(id.ID, userID); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("income not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error deleting income: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted income",
		StatusCode: http.StatusNoContent,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// CreateWealthFundHandler creates a new wealth fund in the database.
//
// @Summary Create a wealth fund
// @Description Create a new wealth fund.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param wealthFund body WealthFundRequest true "Wealth fund object"
// @Success 201 {object} jsonresponse.IdResponse "Wealth fund created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating wealth fund"
// @Security JWT
// @Router /analytics/wealth_fund [post]
func (h *MyHandler) CreateWealthFundHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new wealth fund...")

	// Decode the request payload
	var wealthFundR WealthFundRequest
	if err := json.NewDecoder(r.Body).Decode(&wealthFundR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	wealthFund := wealthFundR.WealthFund

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the wealth fund
	wealthFund.UserID = userID

	// Create a new wealth fund in the database
	wealthFundID, err := h.m.WealthFunds.Create(&wealthFund)
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

// UpdateWealthFundHandler handles the update of an existing expense.
//
// @Summary Update the wealth fund
// @Description Update an existing wealth fund. There is no need to fill user_id field.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param ConnectedAccount body WealthFundRequest true "wealth fund object"
// @Success 200 {object} jsonresponse.SuccessResponse "wealth fund updated successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 404 {object} jsonresponse.ErrorResponse "wealth fund not found"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error updating wealth fund"
// @Security JWT
// @Router /analytics/wealth_fund [put]
func (h *MyHandler) UpdateWealthFundHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Updating wealth fund...")

	// Decode the request body
	var wealthFundR WealthFundRequest
	if err := json.NewDecoder(r.Body).Decode(&wealthFundR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	wealthFund := wealthFundR.WealthFund

	// Check if user is authenticated
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}
	wealthFund.UserID = userID

	// Attempt to update the account
	if err := h.m.WealthFunds.Update(&wealthFund); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("wealth fund not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error updating wealth fund: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	response := jsonresponse.SuccessResponse{
		Message:    "Wealth fund updated successfully",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// DeleteWealthFundHandler handles the deletion of an existing wealth fund.
//
// @Summary Delete the wealth fund
// @Description Delete the existing wealth fund.
// @Tags Analytics
// @Param ConnectedAccount body jsonresponse.IdRequest true "wealth fund id"
// @Success 204 {object} jsonresponse.SuccessResponse "wealth fund deleted successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error deleting wealth fund"
// @Security JWT
// @Router /analytics/wealth_fund [delete]
func (h *MyHandler) DeleteWealthFundHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting wealth fund...")

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

	if err := h.m.WealthFunds.Delete(id.ID, userID); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("wealth fund not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error deleting wealth fund: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted wealth fund",
		StatusCode: http.StatusNoContent,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
