package v1

import (
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/models"
	"go.uber.org/zap"
	"net/http"
)

// CreateExpenseCategoryHandler creates a new expense category in the database.
//
// @Summary CreateExpenseCategoryHandler an expense category
// @Description Creates a new expense category in the database and returns its ID.
// @Tags	Categories
// @Accept 	json
// @Produce json
// @Param 	category body models.ExpenseCategory true "Expense category object"
// @Success 201 {string} string "Expense category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating expense category"
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

	response := map[string]interface{}{
		"message":           "Expense category created successfully",
		"created_object_id": expenseCategoryID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Expense category created successfully", zap.Int64("expenseCategoryID", expenseCategoryID))
}

// CreateIncomeCategoryHandler creates a new income category in the database.
//
// @Summary Create an income category
// @Description Create a new income category.
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.IncomeCategory true "Income category object"
// @Success 201 {string} string "Income category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating income category"
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

	response := map[string]interface{}{
		"message":           "Income category created successfully",
		"created_object_id": incomeCategoryID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Income category created successfully", zap.Int64("incomeCategoryID", incomeCategoryID))
}

// CreateInvestmentCategoryHandler creates a new investment category in the database.
//
// @Summary Create an investment category
// @Description Create a new investment category.
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.InvestmentCategory true "Investment category object"
// @Success 201 {string} string "Investment category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating investment category"
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

	response := map[string]interface{}{
		"message":           "Investment category created successfully",
		"created_object_id": investmentCategoryID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Investment category created successfully", zap.Int64("investmentCategoryID", investmentCategoryID))
}
