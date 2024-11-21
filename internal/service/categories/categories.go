//go:build !exclude_swagger
// +build !exclude_swagger

// Package categories provides functionality related to user analytics, tracking, and additional information.
package categories

import (
	//"encoding/json"

	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"math"
	"time"

	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/repository"
	"github.com/wachrusz/Back-End-API/internal/service/currency"

	"log"
)

type Service struct {
	repo          *mydb.Database
	curr          *currency.Service
	exchangeRates map[string]currency.Valute
}

func NewService(db *mydb.Database, currencyService *currency.Service) *Service {
	return &Service{
		repo:          db,
		curr:          currencyService,
		exchangeRates: make(map[string]currency.Valute),
	}
}

// Analytics represents the structure for analytics data, including income, expense, and wealth fund information.
type Analytics struct {
	Income     []models.Income         `json:"income"`
	Expense    []models.Expense        `json:"expense"`
	WealthFund []repository.WealthFund `json:"wealth_fund"`
}

// Tracker represents the structure for tracking data, including tracking state and goals.
type Tracker struct {
	TrackingState repository.TrackingState `json:"tracking_state"`
	Goal          []models.Goal            `json:"goal"`
	FinHealth     repository.FinHealth     `json:"fin_health"`
}

// More represents additional user information, including app and settings details.
type More struct {
	App      repository.App      `json:"app"`
	Settings repository.Settings `json:"settings"`
}

func round(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return math.Round(num*output) / output
}

func (s *Service) convertCurrency(amount float64, fromCurrencyCode string, toCurrencyCode string) float64 {
	// TODO: refactor: unsafe map exchangeRates
	if fromCurrencyCode == "" || toCurrencyCode == "" {
		return round(amount, 2)
	}

	if fromCurrencyCode == toCurrencyCode {
		return round(amount, 2)
	}

	if len(s.exchangeRates) == 0 {
		s.exchangeRates = s.curr.CurrentCurrencyData.Valute
	}

	switch {
	case toCurrencyCode == "RUB":
		// To RUB
		rate, ok := s.exchangeRates[fromCurrencyCode]
		if !ok {
			log.Printf("Couldn't find exchange rate for %v", fromCurrencyCode)
			return amount
		}
		return round(amount*(rate.Value/float64(rate.Nominal)), 2)

	case fromCurrencyCode == "RUB":
		// From RUB
		rate, ok := s.exchangeRates[toCurrencyCode]
		if !ok {
			log.Printf("Couldn't find exchange rate for %v", toCurrencyCode)
			return amount
		}
		return round(amount/(rate.Value/float64(rate.Nominal)), 2)

	default:
		// Convert non RUB valutes
		rubleRateFrom, ok := s.exchangeRates[fromCurrencyCode]
		if !ok {
			log.Printf("Couldn't find exchange rate for %v", fromCurrencyCode)
			return amount
		}

		rubleRateTo, ok := s.exchangeRates[toCurrencyCode]
		if !ok {
			log.Printf("Couldn't find exchange rate for %v", toCurrencyCode)
			return amount
		}

		return round((amount*(rubleRateFrom.Value/float64(rubleRateFrom.Nominal)))/(rubleRateTo.Value/float64(rubleRateTo.Nominal)), 2)
	}
}

func (s *Service) GetAnalyticsFromDB(userID, currencyCode, limitStr, offsetStr, startDateStr, endDateStr string) (*Analytics, error) {
	// TODO: refactor maybe? too complicated
	if startDateStr == "" {
		startDateStr = time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	}
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	queryIncome := "SELECT id, amount, date, planned, category, sender, connected_account, currency_code FROM income WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC LIMIT $4 OFFSET $5;"
	rowsIncome, err := s.repo.Query(queryIncome, userID, startDateStr, endDateStr, limitStr, offsetStr)
	if err != nil {
		return nil, err
	}
	defer rowsIncome.Close()

	var incomeList []models.Income
	for rowsIncome.Next() {
		var income models.Income
		if err := rowsIncome.Scan(&income.ID, &income.Amount, &income.Date, &income.Planned, &income.CategoryID, &income.Sender, &income.BankAccount, &income.Currency); err != nil {
			return nil, err
		}
		income.UserID = userID
		if income.Currency != currencyCode && currencyCode != "" {
			income.Amount = s.convertCurrency(income.Amount, income.Currency, currencyCode)
		}
		incomeList = append(incomeList, income)
	}

	queryExpense := "SELECT id, amount, date, planned, category, sent_to, connected_account, currency_code FROM expense WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC LIMIT $4 OFFSET $5;"
	rowsExpense, err := s.repo.Query(queryExpense, userID, startDateStr, endDateStr, limitStr, offsetStr)
	if err != nil {
		return nil, err
	}
	defer rowsExpense.Close()

	var expenseList []models.Expense
	for rowsExpense.Next() {
		var expense models.Expense
		if err := rowsExpense.Scan(&expense.ID, &expense.Amount, &expense.Date, &expense.Planned, &expense.CategoryID, &expense.SentTo, &expense.BankAccount, &expense.Currency); err != nil {
			return nil, err
		}
		expense.UserID = userID
		if expense.Currency != currencyCode && currencyCode != "" {
			expense.Amount = s.convertCurrency(expense.Amount, expense.Currency, currencyCode)
		}
		expenseList = append(expenseList, expense)
	}

	queryWealthFund := "SELECT id, amount, date, planned, currency_code, connected_account, user_id, category_id FROM wealth_fund WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date DESC LIMIT $4 OFFSET $5;"
	rowsWealthFund, err := s.repo.Query(queryWealthFund, userID, startDateStr, endDateStr, limitStr, offsetStr)
	if err != nil {
		return nil, err
	}
	defer rowsWealthFund.Close()

	var wealthFundList []repository.WealthFund
	for rowsWealthFund.Next() {
		var wealthFund repository.WealthFund
		if err := rowsWealthFund.Scan(&wealthFund.ID, &wealthFund.Amount, &wealthFund.Date, &wealthFund.PlannedStatus, &wealthFund.Currency, &wealthFund.ConnectedAccount, &wealthFund.UserID, &wealthFund.CategoryID); err != nil {
			return nil, err
		}
		if wealthFund.Currency != currencyCode && currencyCode != "" {
			wealthFund.Amount = s.convertCurrency(wealthFund.Amount, wealthFund.Currency, currencyCode)
		}
		wealthFundList = append(wealthFundList, wealthFund)
	}

	analytics := &Analytics{
		Income:     incomeList,
		Expense:    expenseList,
		WealthFund: wealthFundList,
	}

	return analytics, nil
}

func (s *Service) GetTrackerFromDB(userID, currencyCode, limitStr, offsetStr string) (*Tracker, error) {
	queryGoal := "SELECT id, goal, need, currency, current_state, start_date, end_date FROM goal WHERE user_id = $1 LIMIT $2 OFFSET $3;"
	rowsGoal, err := s.repo.Query(queryGoal, userID, limitStr, offsetStr)
	if err != nil {
		return nil, err
	}
	defer rowsGoal.Close()

	var goalList []models.Goal
	for rowsGoal.Next() {
		var goal models.Goal
		if err := rowsGoal.Scan(&goal.ID, &goal.Goal, &goal.Need, &goal.Currency, &goal.CurrentState, &goal.StartDate, &goal.EndDate); err != nil {
			return nil, err
		}
		goal.UserID = userID
		goal.Need = s.convertCurrency(goal.Need, goal.Currency, currencyCode)
		goalList = append(goalList, goal)
	}
	trackingState := &repository.TrackingState{
		State:  s.getTotalState(userID, currencyCode),
		UserID: userID,
	}

	tracker := &Tracker{
		TrackingState: *trackingState,
		Goal:          goalList,
	}

	return tracker, nil
}

func (s *Service) getTotalState(userID string, convertionCode string) float64 {
	var state float64
	query := `
	WITH all_transactions AS (
		SELECT                                     
			income.id,
			CASE                                                    
				WHEN income.currency_code = 'RUB' THEN income.amount
				ELSE income.amount * COALESCE((SELECT rate_to_ruble FROM exchange_rates WHERE currency_code = income.currency_code), 1)
			END AS converted_amount
		FROM
			income
		WHERE
			income.user_id = $1
		UNION ALL
		SELECT
			expense.id,
			CASE
				WHEN expense.currency_code = 'RUB' THEN -expense.amount
				ELSE -expense.amount * COALESCE((SELECT rate_to_ruble FROM exchange_rates WHERE currency_code = expense.currency_code), 1)
			END AS converted_amount                                                                                                      
		FROM                       
			expense
		WHERE
			expense.user_id = $1
	)                             
	SELECT
		COALESCE(SUM(converted_amount), 0) AS total_balance_in_rubles
	FROM                                                                                                             
		all_transactions;
	
	`
	err := s.repo.QueryRow(query, userID).Scan(&state)
	if err != nil {
		log.Println(err)
		return 0
	}
	return s.convertCurrency(state, "RUB", convertionCode)
}

func (s *Service) GetUserInfoFromDB(userID string) (string, string, error) {
	query := "SELECT surname, name FROM users WHERE id = $1"
	var surname, name string

	row := s.repo.QueryRow(query, userID)
	err := row.Scan(&surname, &name)
	if err != nil {
		return "", "", err
	}

	return surname, name, nil
}

func (s *Service) GetMoreFromDB(userID string) (*More, error) {
	var more More

	subs, err := s.GetSubscriptionFromDB(userID)
	if err != nil {
		log.Println("Error getting Subs from DB:", err)
		return nil, err
	}

	var settings repository.Settings

	app, err := s.GetAppFromDB(userID)
	if err != nil {
	}

	settings.Subscriptions = *subs

	more.App = *app
	more.Settings = settings

	return &more, nil
}

func (s *Service) GetAppFromDB(userID string) (*repository.App, error) {
	connectedAccounts, err := s.GetConnectedAccountsFromDB(userID)
	if err != nil {
		return nil, err
	}

	categorySettings, err := s.GetCategorySettingsFromDB(userID)
	if err != nil {
		return nil, err
	}

	app := &repository.App{
		ConnectedAccounts: connectedAccounts,
		CategorySettings:  *categorySettings,
		//OperationArchive:  operationArchive,
	}

	return app, nil
}

func (s *Service) GetSubscriptionFromDB(userID string) (*repository.Subscription, error) {
	var subscription repository.Subscription

	query := "SELECT id, user_id, start_date, end_date, is_active FROM subscriptions WHERE user_id = $1"
	row := s.repo.QueryRow(query, userID)

	err := row.Scan(&subscription.ID, &subscription.UserID, &subscription.StartDate, &subscription.EndDate, &subscription.IsActive)
	if err != nil {
		return &repository.Subscription{}, nil
	}

	return &subscription, nil
}

func (s *Service) GetConnectedAccountsFromDB(userID string) ([]models.ConnectedAccount, error) {
	var connectedAccounts []models.ConnectedAccount

	// Запрос к базе данных для выбора подключенных аккаунтов по идентификатору пользователя.
	query := `
		SELECT id, user_id, bank_id, account_number, account_type
		FROM connected_accounts
		WHERE user_id = $1;
	`

	rows, err := s.repo.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var connectedAccount models.ConnectedAccount
		err := rows.Scan(
			&connectedAccount.ID,
			&connectedAccount.UserID,
			&connectedAccount.BankID,
			&connectedAccount.AccountNumber,
			&connectedAccount.AccountType,
		)
		if err != nil {
			return nil, err
		}

		connectedAccounts = append(connectedAccounts, connectedAccount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return connectedAccounts, nil
}

func (s *Service) GetCategorySettingsFromDB(userID string) (*repository.CategorySettings, error) {
	var categorySettings repository.CategorySettings

	// Запрос для получения конфигурации доходов
	queryIncome := "SELECT id, name, icon, is_fixed, user_id FROM income_categories WHERE user_id = $1"
	rowsIncome, err := mydb.GlobalDB.Query(queryIncome, userID)
	if err != nil {
		log.Println("Error getting income category configuration from DB:", err)
		return nil, err
	}
	defer rowsIncome.Close()

	for rowsIncome.Next() {
		var config repository.IncomeCategory
		err := rowsIncome.Scan(&config.ID, &config.Name, &config.Icon, &config.IsConstant, &config.UserID)
		if err != nil {
			log.Println("Error scanning income category configuration:", err)
			return nil, err
		}
		categorySettings.IncomeCategories = append(categorySettings.IncomeCategories, config)
	}

	// Запрос для получения конфигурации расходов
	queryExpense := "SELECT id, name, icon, is_fixed, user_id FROM expense_categories WHERE user_id = $1"
	rowsExpense, err := s.repo.Query(queryExpense, userID)
	if err != nil {
		log.Println("Error getting expense category configuration from DB:", err)
		return nil, err
	}
	defer rowsExpense.Close()

	for rowsExpense.Next() {
		var config repository.ExpenseCategory
		err := rowsExpense.Scan(&config.ID, &config.Name, &config.Icon, &config.IsConstant, &config.UserID)
		if err != nil {
			log.Println("Error scanning expense category configuration:", err)
			return nil, err
		}
		categorySettings.ExpenseCategories = append(categorySettings.ExpenseCategories, config)
	}

	queryInvestment := "SELECT id, name, icon, is_fixed, user_id FROM investment_categories WHERE user_id = $1"
	rowsInvestment, err := s.repo.Query(queryInvestment, userID)
	if err != nil {
		log.Println("Error getting investment category configuration from DB:", err)
		return nil, err
	}
	defer rowsInvestment.Close()

	for rowsInvestment.Next() {
		var config repository.InvestmentCategory
		err := rowsInvestment.Scan(&config.ID, &config.Name, &config.Icon, &config.IsConstant, &config.UserID)
		if err != nil {
			log.Println("Error scanning investment category configuration:", err)
			return nil, err
		}
		categorySettings.InvestmentCategories = append(categorySettings.InvestmentCategories, config)
	}

	// Проверка, что были получены данные
	if len(categorySettings.ExpenseCategories) == 0 && len(categorySettings.IncomeCategories) == 0 && len(categorySettings.InvestmentCategories) == 0 {
		return &repository.CategorySettings{}, nil
	}

	return &categorySettings, nil
}

func (s *Service) GetOperationArchiveFromDB(userID, limit, offset string) ([]repository.Operation, error) {
	var operations []repository.Operation

	query := `
		SELECT id, description, amount, date, category, operation_type
		FROM operations
		WHERE user_id = $1
		ORDER BY date DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := s.repo.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var operation repository.Operation
		err := rows.Scan(
			&operation.ID,
			&operation.Description,
			&operation.Amount,
			&operation.Date,
			&operation.Category,
			&operation.Type,
		)
		if err != nil {
			return nil, err
		}

		operations = append(operations, operation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return operations, nil
}

type Categories interface {
	GetAnalyticsFromDB(userID, currencyCode, limitStr, offsetStr, startDateStr, endDateStr string) (*Analytics, error)
	GetTrackerFromDB(userID, currencyCode, limitStr, offsetStr string) (*Tracker, error)
	GetUserInfoFromDB(userID string) (string, string, error)
	GetMoreFromDB(userID string) (*More, error)
	GetAppFromDB(userID string) (*repository.App, error)
	GetSubscriptionFromDB(userID string) (*repository.Subscription, error)
	GetConnectedAccountsFromDB(userID string) ([]models.ConnectedAccount, error)
	GetCategorySettingsFromDB(userID string) (*repository.CategorySettings, error)
	GetOperationArchiveFromDB(userID, limit, offset string) ([]repository.Operation, error)
}
