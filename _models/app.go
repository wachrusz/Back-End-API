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

func CreateIncomeCategory(category *IncomeCategory) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO income_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4)",
		category.Name, category.Icon, category.IsConstant, category.UserID)
	if err != nil {
		log.Println("Error creating income:", err)
		return err
	}
	return nil
}

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

func CreateInvestmentCategory(category *InvestmentCategory) error {
	_, err := mydb.GlobalDB.Exec("INSERT INTO investment_categories (name, icon, is_fixed, user_id) VALUES ($1, $2, $3, $4)",
		category.Name, category.Icon, category.IsConstant, category.UserID)
	if err != nil {
		log.Println("Error creating investment:", err)
		return err
	}
	return nil
}
