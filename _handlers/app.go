package handlers

import (
	models "backEndAPI/_models"
	"encoding/json"
	"net/http"
)

// CreateExpenseCategoryHandler обрабатывает запрос на создание категории расходов.
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

// CreateIncomeCategoryHandler обрабатывает запрос на создание категории доходов.
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
