//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "backEndAPI/_mydatabase"
	"log"
)

// App представляет собой информацию о приложении пользователя.
type App struct {
	ConnectedAccounts []ConnectedAccount `json:"connected_accounts"`
	CategorySettings  CategorySettings   `json:"category_settings"`
	OperationArchive  []Operation        `json:"operation_archive"`
}

// ConnectedAccount представляет собой информацию о подключенном счете.
type ConnectedAccount struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

// CategorySettings представляет собой настройки категорий.
type CategorySettings struct {
	ExpenseCategories    []ExpenseCategory    `json:"expense_categories"`
	IncomeCategories     []IncomeCategory     `json:"income_categories"`
	InvestmentCategories []InvestmentCategory `json:"investment_category"`
}

// Operation представляет собой информацию об операции.
type Operation struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
}

// ExpenseCategory представляет собой информацию о категории расходов.
type ExpenseCategory struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	IsConstant bool   `json:"is_constant"`
	UserID     string `json:"user_id"`
}

// IncomeCategory представляет собой информацию о категории доходов.
type IncomeCategory struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	IsConstant bool   `json:"is_constant"`
	UserID     string `json:"user_id"`
}

type InvestmentCategory struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	IsConstant bool   `json:"is_constant"`
	UserID     string `json:"user_id"`
}

// @Summary Create income category
// @Description Create a new income category.
// @Tags Income
// @Accept json
// @Produce json
// @Param category body IncomeCategory true "Income category details"
// @Success 201 {string} string "Income category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating income category"
// @Router /models/income/category [post]
func CreateIncomeCategory(category *IncomeCategory) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO income_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4)",
		category.Name, category.Icon, category.IsConstant, category.UserID)
	if err != nil {
		log.Println("Error creating income:", err)
		return err
	}
	return nil
}

// @Summary Create expense category
// @Description Create a new expense category.
// @Tags Expense
// @Accept json
// @Produce json
// @Param category body ExpenseCategory true "Expense category details"
// @Success 201 {string} string "Expense category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating expense category"
// @Router /models/expense/category [post]
func CreateExpenseCategory(category *ExpenseCategory) error {
	log.Println("category: ", category)
	_, err := mydb.GlobalDB.Exec("INSERT INTO expense_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4)",
		category.Name, category.Icon, category.IsConstant, category.UserID)
	if err != nil {
		log.Println("Error creating expense:", err)
		return err
	}
	return nil
}

// @Summary Create investment category
// @Description Create a new investment category.
// @Tags Investment
// @Accept json
// @Produce json
// @Param category body InvestmentCategory true "Investment category details"
// @Success 201 {string} string "Investment category created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error creating investment category"
// @Router /models/investment/category [post]
func CreateInvestmentCategory(category *InvestmentCategory) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO investment_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4)",
		category.Name, category.Icon, category.IsConstant, category.UserID)
	if err != nil {
		log.Println("Error creating investment:", err)
		return err
	}
	return nil
}
