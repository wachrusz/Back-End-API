//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

import (
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"log"
)

// App представляет собой информацию о приложении пользователя.
type App struct {
	ConnectedAccounts []ConnectedAccount `json:"connected_accounts"`
	CategorySettings  CategorySettings   `json:"category_settings"`
	//OperationArchive  []Operation        `json:"operation_archive"` //*Deleted from APP
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

func CreateIncomeCategory(category *IncomeCategory) (int64, error) {
	var incomeCategoryID int64
	err := mydb.GlobalDB.QueryRow("INSERT INTO income_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		category.Name, category.Icon, category.IsConstant, category.UserID).Scan(&incomeCategoryID)
	if err != nil {
		log.Println("Error creating income:", err)
		return 0, err
	}
	return incomeCategoryID, nil
}

func CreateExpenseCategory(category *ExpenseCategory) (int64, error) {
	var expenseCategoryID int64
	err := mydb.GlobalDB.QueryRow("INSERT INTO expense_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		category.Name, category.Icon, category.IsConstant, category.UserID).Scan(&expenseCategoryID)
	if err != nil {
		log.Println("Error creating expense:", err)
		return 0, err
	}
	return expenseCategoryID, nil
}

func CreateInvestmentCategory(category *InvestmentCategory) (int64, error) {
	var investmentCategoryID int64
	err := mydb.GlobalDB.QueryRow("INSERT INTO investment_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
		category.Name, category.Icon, category.IsConstant, category.UserID).Scan(&investmentCategoryID)
	if err != nil {
		log.Println("Error creating investment:", err)
		return 0, err
	}
	return investmentCategoryID, nil
}
