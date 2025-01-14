package fin_health

import "math"

// ExpensePropensity считает среднюю склонность к потреблению
// Формула: Суммарные расходы за месяц/располагаемый доход за месяц
// Формула преобразования: min{100*(1.2-propensity_to_expend); 50}
func (s *Service) ExpensePropensity(userID string) (float64, error) {
	q := `
	WITH current_month_expenses AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM expense_in_rubles
		WHERE
			user_id = $1 AND
			planned = false AND
			date >= NOW() - INTERVAL '30 days'
	), 
	current_month_incomes AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM income_in_rubles
		WHERE
			user_id = $1 AND
			planned = false AND
			date >= NOW() - INTERVAL '30 days'
	) 
	SELECT 
	    CASE
			WHEN (incomes.total = 0) OR (incomes.total - expenses.total = 0) THEN 0
			ELSE expenses.total / (incomes.total - expenses.total)
		END AS propensity,
	    incomes.total AS incomes,
	    expenses.total AS expenses
	FROM 
		current_month_expenses expenses, 
		current_month_incomes incomes
	`

	var propensity, incomes, expenses float64
	err := s.repo.QueryRow(q, userID).Scan(&propensity, &incomes, &expenses)
	if err != nil {
		return 0, err
	}

	if incomes == 0 && expenses != 0 {
		return 0, nil
	}

	result := math.Min(100*(1.2-propensity), 50)
	return result, nil
}

// ExpenditureDelta считает изменение расходов по сравнению со среднемесячными расходами
// Формула: (суммарные расходы за данный месяц - средние ежемесячные расходы за последние 3 месяца)/средние ежемесячные расходы за последние 3 месяца *100
// Формула преобразования: min{2.5*(15 - expenditure_delta); 50}
func (s *Service) ExpenditureDelta(userID string) (float64, error) {
	q := `
	WITH current_month_expenses AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM expense_in_rubles
		WHERE
			user_id = $1 AND
			planned = false AND
			date >= NOW() - INTERVAL '30 days'
	),
	average_monthly_expenses AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) / 3 AS three_months
		FROM expense_in_rubles
		WHERE
			user_id = $1 AND
			planned = false AND
			date >= NOW() - INTERVAL '90 days'
	)
	SELECT
	    CASE
	        WHEN average.three_months = 0 THEN 0
            ELSE (current.total - average.three_months) / average.three_months * 100
		END AS delta,
        average.three_months AS avg
	FROM 
		current_month_expenses current, 
		average_monthly_expenses average;
	`

	var delta, avg float64
	err := s.repo.QueryRow(q, userID).Scan(&delta, &avg)
	if err != nil {
		return 0, err
	}

	if avg == 0 {
		return 0, nil
	}

	result := math.Min(2.5*(15.0-delta), 50)
	return result, nil
}
