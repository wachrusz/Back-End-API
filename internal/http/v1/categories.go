package v1

import (
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/models"
	"github.com/wachrusz/Back-End-API/pkg/json_response"
	"go.uber.org/zap"
	"net/http"
)

// CreateExpenseCategoryHandler creates a new expense category in the database.
//
// @Summary CreateExpenseCategoryHandler an expense category
// @Description Creates a new expense category in the database and returns its ID.
// @Tags	App
// @Accept 	json
// @Produce json
// @Param 	category body models.ExpenseCategory true "Expense category object"
// @Success 201 {object} jsonresponse.IdResponse "Expense category created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating expense category"
// @Security JWT
// @Router /app/category/expense [post]
func (h *MyHandler) CreateExpenseCategoryHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating category...")
	var category models.ExpenseCategory

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	expenseCategoryID, err := models.CreateExpenseCategory(&category)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating expense category: %v", err), http.StatusInternalServerError)
		return
	}

	response := jsonresponse.IdResponse{
		Message:    "Expense category created successfully",
		Id:         expenseCategoryID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Expense category created successfully", zap.Int64("expenseCategoryID", expenseCategoryID))
}

// CreateIncomeCategoryHandler creates a new income category in the database.
//
// @Summary Create an income category
// @Description Create a new income category.
// @Tags App
// @Accept json
// @Produce json
// @Param category body models.IncomeCategory true "Income category object"
// @Success 201 {object} jsonresponse.IdResponse "Income category created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating income category"
// @Security JWT
// @Router /app/category/income [post]
func (h *MyHandler) CreateIncomeCategoryHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating income category...")

	var category models.IncomeCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	incomeCategoryID, err := models.CreateIncomeCategory(&category)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating income category: %v", err), http.StatusInternalServerError)
		return
	}

	response := jsonresponse.IdResponse{
		Message:    "Income category created successfully",
		Id:         incomeCategoryID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Income category created successfully", zap.Int64("incomeCategoryID", incomeCategoryID))
}

// CreateInvestmentCategoryHandler creates a new investment category in the database.
//
// @Summary Create an investment category
// @Description Create a new investment category.
// @Tags App
// @Accept json
// @Produce json
// @Param category body models.InvestmentCategory true "Investment category object"
// @Success 201 {object} jsonresponse.IdResponse "Investment category created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error creating investment category"
// @Security JWT
// @Router /app/category/investment [post]
func (h *MyHandler) CreateInvestmentCategoryHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Creating investment category...")

	var category models.InvestmentCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	investmentCategoryID, err := models.CreateInvestmentCategory(&category)
	if err != nil {
		h.errResp(w, fmt.Errorf("error creating investment category: %v", err), http.StatusInternalServerError)
		return
	}

	response := jsonresponse.IdResponse{
		Message:    "Investment category created successfully",
		Id:         investmentCategoryID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Investment category created successfully", zap.Int64("investmentCategoryID", investmentCategoryID))
}
