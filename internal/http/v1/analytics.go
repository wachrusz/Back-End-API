package v1

import (
	"encoding/json"
	"fmt"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"go.uber.org/zap"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/repository"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
)

// CreateExpenseHandler creates a new expense record in the database.
//
// @Summary Create an expense
// @Description Create a new expense record.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param expense body repository.Expense true "Expense object"
// @Success 201 {object} jsonresponse.IdResponse "Successfully created an expense"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating expense"
// @Security JWT
// @Router /expenses [post]
func (h *MyHandler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating a new expense...")

	// Decode the request payload
	var expense repository.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Extract the user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Assign user ID to the expense
	expense.UserID = userID

	// Create a new expense in the database
	expenseID, err := repository.CreateExpense(&expense)
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
