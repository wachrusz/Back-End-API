package fin_health

import "math"

// ratioQuery это запрос, для Service.LiquidFundRatio и Service.IlliquidFundRatio. Они имеют схожую структуру.
var ratioQuery string = `
	WITH liquid_active AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total_liquid
		FROM wealth_fund_in_rubles
		WHERE
			user_id = $1  AND
			is_liquid = $2 AND 
			planned = '0' AND
			date >= NOW() - INTERVAL '1 year'
	),
	average_yearly_expense AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) / 12 AS avg_last_year
		FROM expense_in_rubles
		WHERE
			user_id = $1 AND
			planned = '0' AND
			date >= NOW() - INTERVAL '1 year'
	)
	SELECT 
		la.total_liquid / ay.avg_last_year AS ratio
	FROM 
		liquid_active la, 
		average_yearly_expense ay;`

// LiquidFundRatio считает отношение ликвидного фонда благосостояния к средним ежемесячным расходам
// Формула: общая сумма ликвидных активов/средние ежемесячные расходы за последний год
// Формула преобразования: min{liquid_fund_ratio*10; 30}
func (s *Service) LiquidFundRatio(userID string) (float64, error) {
	var ratio float64
	err := s.repo.QueryRow(ratioQuery, userID, liquid).Scan(&ratio)
	if err != nil {
		return 0, err
	}

	result := math.Min(ratio*10, 30)
	return result, nil
}

// IlliquidFundRatio считает отношение неликвидного фонда благосостояния к средним ежемесячным расходам
// Формула: общая сумма неликвидных активов/средние ежемесячные расходы за последний год
// Формула преобразования: min{illiquid_fund_ratio*5/3; 20}
func (s *Service) IlliquidFundRatio(userID string) (float64, error) {
	var ratio float64
	err := s.repo.QueryRow(ratioQuery, userID, illiquid).Scan(&ratio)
	if err != nil {
		return 0, err
	}

	result := math.Min(ratio*10, 30)
	return result, nil
}

// SavingsToIncomeRatio считает отношение отчислений на сбережения относительно дохода за месяц
// Формула: общая сумма отчислений на сбережения/общая сумма доходов за месяц
// Формула преобразования: min{saving_to_income_ratio*150; 30}
func (s *Service) SavingsToIncomeRatio(userID string) (float64, error) {
	q := `
	WITH expense_for_savings AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM expense_in_rubles
		WHERE
			user_id = $1  AND
			type = $2 AND 
			planned = '0' AND
			date >= NOW() - INTERVAL '30 days'
	),
	income_monthly AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM income_in_rubles
		WHERE
			user_id = $1 AND
			planned = '0' AND
			date >= NOW() - INTERVAL '30 days'
	)
	SELECT
	    CASE 
			WHEN income_monthly.total THEN 0
	    	ELSE expense_for_savings.total / income_monthly.total
	    END AS ratio
	FROM 
		expense_for_savings, 
		income_monthly;
	`

	var ratio float64
	err := s.repo.QueryRow(q, userID, saving).Scan(&ratio)
	if err != nil {
		return 0, err
	}

	result := math.Min(ratio*150, 30)
	return result, nil
}

// SavingDelta считает изменение накоплений по сравнению со среднемесячными
// Формула: ((сбереженная сумма за данный месяц - средняя сбереженная сумма за последний год)/средняя сбереженная сумма за последний год) +1
// Формула преобразования: min{(delta-0.8)*50; 20}
func (s *Service) SavingDelta(userID string) (float64, error) {
	q := `
	WITH current_month_savings AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM wealth_fund_in_rubles
		WHERE
			user_id = $1 AND
			type = $2 AND
			planned = '0' AND
			date >= NOW() - INTERVAL '30 days'
	), 
	average_saving_amount_annually AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) / 12 AS avg_amount
		FROM wealth_fund_in_rubles
		WHERE
			user_id = $1 AND
			type = $2 AND
			planned = '0' AND
			date >= NOW() - INTERVAL '1 year'
	) 
	SELECT 
	    CASE
			WHEN av.avg_amount = 0 THEN 0
			ELSE (cs.total - av.avg_amount) / av.avg_amount + 1
		END AS delta 
	FROM 
		current_month_savings cs, 
		average_saving_amount_annually av
	`

	var delta float64
	err := s.repo.QueryRow(q, userID, saving).Scan(&delta)
	if err != nil {
		return 0, err
	}

	result := math.Min((delta-0.8)*50, 20)
	return result, nil
}
