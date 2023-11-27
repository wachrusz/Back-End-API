//go:build !exclude_swagger
// +build !exclude_swagger

// Package categories provides functionality related to user analytics, tracking, and additional information.
package categories

import (
	//"encoding/json"

	logger "backEndAPI/_logger"
	models "backEndAPI/_models"
	mydb "backEndAPI/_mydatabase"
	"log"
)

// Analytics represents the structure for analytics data, including income, expense, and wealth fund information.
type Analytics struct {
	Income     []models.Income
	Expense    []models.Expense
	WealthFund []models.WealthFund
}

// Tracker represents the structure for tracking data, including tracking state and goals.
type Tracker struct {
	TrackingState models.TrackingState
	Goal          []models.Goal
	//! TODO: FinHealth models.FinHealth
}

// More represents additional user information, including app and settings details.
type More struct {
	App      models.App
	Settings models.Settings
}

// @Summary Get user analytics from the database
// @Description Get income, expense, and wealth fund data for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} Analytics "User analytics data"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/analytics/{userID} [get]
func GetAnalyticsFromDB(userID string) (*Analytics, error) {

	queryIncome := "SELECT id, amount, date, planned FROM income WHERE user_id = $1"
	rowsIncome, err := mydb.GlobalDB.Query(queryIncome, userID)
	if err != nil {
		return nil, err
	}
	defer rowsIncome.Close()

	var incomeList []models.Income
	for rowsIncome.Next() {
		var income models.Income
		if err := rowsIncome.Scan(&income.ID, &income.Amount, &income.Date, &income.Planned); err != nil {
			return nil, err
		}
		income.UserID = userID
		incomeList = append(incomeList, income)
	}

	queryExpense := "SELECT id, amount, date, planned FROM expense WHERE user_id = $1"
	rowsExpense, err := mydb.GlobalDB.Query(queryExpense, userID)
	if err != nil {
		return nil, err
	}
	defer rowsExpense.Close()

	var expenseList []models.Expense
	for rowsExpense.Next() {
		var expense models.Expense
		if err := rowsExpense.Scan(&expense.ID, &expense.Amount, &expense.Date, &expense.Planned); err != nil {
			return nil, err
		}
		expense.UserID = userID
		expenseList = append(expenseList, expense)
	}

	queryWealthFund := "SELECT id, amount, date FROM wealth_fund WHERE user_id = $1"
	rowsWealthFund, err := mydb.GlobalDB.Query(queryWealthFund, userID)
	if err != nil {
		return nil, err
	}
	defer rowsWealthFund.Close()

	var wealthFundList []models.WealthFund
	for rowsWealthFund.Next() {
		var wealthFund models.WealthFund
		if err := rowsWealthFund.Scan(&wealthFund.ID, &wealthFund.Amount, &wealthFund.Date); err != nil {
			return nil, err
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

// @Summary Get user tracker data from the database
// @Description Get tracking state and goal data for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Param analytics body Analytics true "User analytics data"
// @Success 200 {object} Tracker "User tracker data"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/tracker/{userID} [get]
func GetTrackerFromDB(userID string, analytics *Analytics) (*Tracker, error) {
	queryGoal := "SELECT id, goal, need, current_state FROM goal WHERE user_id = $1"
	rowsGoal, err := mydb.GlobalDB.Query(queryGoal, userID)
	if err != nil {
		logger.ErrorLogger.Print("Error getting Goal From DB: (userID, error) ", userID, err)
		return nil, err
	}
	defer rowsGoal.Close()

	var goalList []models.Goal
	for rowsGoal.Next() {
		var goal models.Goal
		if err := rowsGoal.Scan(&goal.ID, &goal.Goal, &goal.Need, &goal.CurrentState); err != nil {
			return nil, err
		}
		goal.UserID = userID
		goalList = append(goalList, goal)
	}
	trackingState := &models.TrackingState{
		State:  getTotalState(analytics),
		UserID: userID,
	}

	tracker := &Tracker{
		TrackingState: *trackingState,
		Goal:          goalList,
	}

	return tracker, nil
}

func getTotalState(analytics *Analytics) float64 {
	var Sum float64
	for _, income := range analytics.Income {
		Sum += income.Amount
	}
	for _, expense := range analytics.Expense {
		Sum -= expense.Amount
	}
	return Sum
}

// @Summary Get user information from the database
// @Description Get username and name information for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} string "User information"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/userinfo/{userID} [get]
func GetUserInfoFromDB(userID string) (string, string, error) {
	query := "SELECT username, name FROM users WHERE id = $1"
	var username, name string

	row := mydb.GlobalDB.QueryRow(query, userID)
	err := row.Scan(&username, &name)
	if err != nil {
		logger.ErrorLogger.Print("Error getting user information from DB: ", err)
		return "", "", err
	}

	return username, name, nil
}

// @Summary Get additional user information from the database
// @Description Get more details, including app and settings information, for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} More "Additional user information"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/more/{userID} [get]
func GetMoreFromDB(userID string) (*More, error) {
	var more More

	subs, err := GetSubscriptionFromDB(userID)
	if err != nil {
		log.Println("Error getting Subs from DB:", err)
		return nil, err
	}

	var settings models.Settings

	app, err := GetAppFromDB(userID)
	if err != nil {
		logger.ErrorLogger.Printf("Error in GetAppFromDB: %v", err)
	}

	settings.Subscriptions = *subs

	more.App = *app
	more.Settings = settings

	return &more, nil
}

// @Summary Get app information from the database
// @Description Get app details, including connected accounts, category settings, and operation archive, for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} models.App "App information"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/app/{userID} [get]
func GetAppFromDB(userID string) (*models.App, error) {
	connectedAccounts, err := GetConnectedAccountsFromDB(userID)
	if err != nil {
		return nil, err
	}

	categorySettings, err := GetCategorySettingsFromDB(userID)
	if err != nil {
		return nil, err
	}

	operationArchive, err := GetOperationArchiveFromDB(userID)
	if err != nil {
		return nil, err
	}

	app := &models.App{
		ConnectedAccounts: connectedAccounts,
		CategorySettings:  *categorySettings,
		OperationArchive:  operationArchive,
	}

	return app, nil
}

// @Summary Get subscription information from the database
// @Description Get subscription details for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} models.Subscription "Subscription information"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/subscription/{userID} [get]
func GetSubscriptionFromDB(userID string) (*models.Subscription, error) {
	var subscription models.Subscription

	query := "SELECT id, user_id, start_date, end_date, is_active FROM subscriptions WHERE user_id = $1"
	row := mydb.GlobalDB.QueryRow(query, userID)

	err := row.Scan(&subscription.ID, &subscription.UserID, &subscription.StartDate, &subscription.EndDate, &subscription.IsActive)
	if err != nil {
		log.Println("Error getting subscription information from DB:", err)
		return &models.Subscription{}, nil
	}

	return &subscription, nil
}

// @Summary Get connected accounts information from the database
// @Description Get connected accounts details for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} []models.ConnectedAccount "Connected accounts information"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/connectedaccounts/{userID} [get]
func GetConnectedAccountsFromDB(userID string) ([]models.ConnectedAccount, error) {
	// Здесь реализация получения связанных счетов из базы данных
	//! TODO
	return []models.ConnectedAccount{}, nil
}

// @Summary Get category settings information from the database
// @Description Get category settings, including income, expense, and investment category details, for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} models.CategorySettings "Category settings information"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/categorysettings/{userID} [get]
func GetCategorySettingsFromDB(userID string) (*models.CategorySettings, error) {
	var categorySettings models.CategorySettings

	// Запрос для получения конфигурации доходов
	queryIncome := "SELECT id, name, icon, is_fixed, user_id FROM income_categories WHERE user_id = $1"
	rowsIncome, err := mydb.GlobalDB.Query(queryIncome, userID)
	if err != nil {
		log.Println("Error getting income category configuration from DB:", err)
		return nil, err
	}
	defer rowsIncome.Close()

	for rowsIncome.Next() {
		var config models.IncomeCategory
		err := rowsIncome.Scan(&config.ID, &config.Name, &config.Icon, &config.IsConstant, &config.UserID)
		if err != nil {
			log.Println("Error scanning income category configuration:", err)
			return nil, err
		}
		categorySettings.IncomeCategories = append(categorySettings.IncomeCategories, config)
	}

	// Запрос для получения конфигурации расходов
	queryExpense := "SELECT id, name, icon, is_fixed, user_id FROM expense_categories WHERE user_id = $1"
	rowsExpense, err := mydb.GlobalDB.Query(queryExpense, userID)
	if err != nil {
		log.Println("Error getting expense category configuration from DB:", err)
		return nil, err
	}
	defer rowsExpense.Close()

	for rowsExpense.Next() {
		var config models.ExpenseCategory
		err := rowsExpense.Scan(&config.ID, &config.Name, &config.Icon, &config.IsConstant, &config.UserID)
		if err != nil {
			log.Println("Error scanning expense category configuration:", err)
			return nil, err
		}
		categorySettings.ExpenseCategories = append(categorySettings.ExpenseCategories, config)
	}

	queryInvestment := "SELECT id, name, icon, is_fixed, user_id FROM investment_categories WHERE user_id = $1"
	rowsInvestment, err := mydb.GlobalDB.Query(queryInvestment, userID)
	if err != nil {
		log.Println("Error getting investment category configuration from DB:", err)
		return nil, err
	}
	defer rowsInvestment.Close()

	for rowsInvestment.Next() {
		var config models.InvestmentCategory
		err := rowsInvestment.Scan(&config.ID, &config.Name, &config.Icon, &config.IsConstant, &config.UserID)
		if err != nil {
			log.Println("Error scanning investment category configuration:", err)
			return nil, err
		}
		categorySettings.InvestmentCategories = append(categorySettings.InvestmentCategories, config)
	}

	// Проверка, что были получены данные
	if len(categorySettings.ExpenseCategories) == 0 && len(categorySettings.IncomeCategories) == 0 && len(categorySettings.InvestmentCategories) == 0 {
		return &models.CategorySettings{}, nil
	}

	return &categorySettings, nil
}

// @Summary Get operation archive information from the database
// @Description Get operation archive details for a specific user from the database.
// @Tags Categories
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} []models.Operation "Operation archive information"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories/operationarchive/{userID} [get]
func GetOperationArchiveFromDB(userID string) ([]models.Operation, error) {
	// Здесь реализация получения архива операций из базы данных
	//! TODO

	return []models.Operation{}, nil
}
