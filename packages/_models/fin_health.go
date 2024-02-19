//go:build !exclude_swagger
// +build !exclude_swagger

// Package models provides basic financial models functionality.
package models

//FinHealth
type FinHealth struct {
	ID              string `json:"id"`
	IncomeScore     int    `json:"income_score"`
	ExpenseScore    int    `json:"expense_score"`
	InvestmentScore int    `json:"investment_score"`
	ObligationScore int    `json:"obligation_score"`
	PlanScore       int    `json:"plan_score"`
	TotalScore      int    `json:"total_score"`
	UserID          string `json:"user_id"`
}

/*
func CreateFinHealth(finHealth *FinHealth) error {
	finHealth.IncomeScore = getIncomeScore(finHealth.UserID)

	_, err1 := mydb.GlobalDB.Exec("INSERT INTO wealth_fund (amount, date, user_id) VALUES ($1, $2, $3)",
		wealthFund.Amount, parsedDate, wealthFund.UserID)
	if err1 != nil {
		log.Println("Error creating wealthFund:", err)
		return err1
	}
	return nil
}

func getIncomeScore(userID string) int {
	query := `
		SELECT id, amount FROM income WHERE user_id = $1;
	`

	rows, err := mydb.Database.Query(query, userID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var totalScore int
	for rows.Next() {
		var income Income
		err := rows.Scan(&income.ID, &income.Amount)
		if err != nil {
			log.Fatal(err)
		}

		incomeScore := int(income.Amount) * 10

		totalScore += incomeScore
	}

	return totalScore
}

func FinHealthCalculations(userID string) (*FinHealth, error) {
	var finHealth FinHealth

	incomeTemp, incomePerc, err := GetMonthlyIncomeIncrease(userID)
	if err != nil {
		return &finHealth, err
	}

	expenseTemp, expensePerc, err := GetMonthlyExpenseIncrease(userID)
	if err != nil {
		return &finHealth, err
	}

	incomeExpenseDiff := incomeTemp - expenseTemp
	if incomeExpenseDiff < 0 {
		incomeExpenseDiff = 0
	}

	return &finHealth, nil
}
*/
