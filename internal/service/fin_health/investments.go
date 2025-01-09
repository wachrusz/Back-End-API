package fin_health

import "math"

// InvestmentsToSavingsRatio считает ежемесячное отчисление на инвестиции относительно ежемесячных отчислений на сбережения
// Формула: сумма отчислений на инвестиции/сумма отчислений на сбережения за последний месяц
// Формула преобразования: min{monthly_investment_to_savings_ratio*40; 20}
func (s *Service) InvestmentsToSavingsRatio(userID string) (float64, error) {
	q := `
	WITH investments AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS monthly
		FROM expense_in_rubles
		WHERE
			user_id = $1 AND
			planned = '0' AND
			type = $2 AND
			date >= NOW() - INTERVAL '30 days'
	),
	savings AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS monthly
		FROM expense_in_rubles
		WHERE
			user_id = $1 AND
			planned = '0' AND
			type = $3 AND
			date >= NOW() - INTERVAL '30 days'
	)
	SELECT 
	    CASE
            WHEN s.monthly = 0 THEN 0
			ELSE i.monthly / s.monthly
	    END AS ratio
	FROM 
		investments i, 
		savings s;
	`

	var ratio float64
	err := s.repo.QueryRow(q, userID, investment, saving).Scan(&ratio)
	if err != nil {
		return 0, err
	}

	result := math.Min(ratio*40, 20)
	return result, nil
}

// InvestmentsToFundRatio считает долю отчислений на инвестиции относительно накоплений
// Формула: общая сумма инвестиций/общая сумма накоплений
// Формула преобразования: min{investment_to_fund_ratio*100; 50}
func (s *Service) InvestmentsToFundRatio(userID string) (float64, error) {
	q := `
	WITH investments AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM wealth_fund_in_rubles
		WHERE
			user_id = $1 AND
			planned = '0' AND
			type = $2
	),
	savings AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM wealth_fund_in_rubles
		WHERE
			user_id = $1 AND
			planned = '0' AND
			type = $3
	)
	SELECT 
	    CASE
            WHEN s.total = 0 THEN 0
			ELSE i.total / s.total
	    END AS ratio
	FROM 
		investments i, 
		savings s;
	`

	var ratio float64
	err := s.repo.QueryRow(q, userID, investment, saving).Scan(&ratio)
	if err != nil {
		return 0, err
	}

	result := math.Min(ratio*100, 50)
	return result, nil
}
