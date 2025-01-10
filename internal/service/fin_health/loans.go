package fin_health

import "math"

// LoansToAssetsRatio считает отношение общей суммы обязательств к фонду благосостояния
// Формула: общая сумма задолженностей человека/суммарный фонд благосостояния
// Формула преобразования: min{90(0.5-loans_to_assets); 45}
func (s *Service) LoansToAssetsRatio(userID string) (float64, error) {
	q := `
	WITH loans AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total,
		       COALESCE(COUNT(amount_in_rubles), 0) AS count
		FROM wealth_fund_in_rubles
		WHERE
			user_id = $1  AND
			type = $2 AND 
			planned = false	
	),
	fund AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM wealth_fund_in_rubles
		WHERE
			user_id = $1 AND
			planned = false
	)
	SELECT 
	    CASE
	        WHEN fund.total = 0 THEN 0
			ELSE loans.total / fund.total
	    END AS ratio, 
	    loans.count AS count
	FROM loans, fund;
	`

	var ratio float64
	var count int
	err := s.repo.QueryRow(q, userID, loan).Scan(&ratio, &count)
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 100, nil
	}

	result := math.Min(90*(0.5-ratio), 45)
	return result, nil
}

// LoansPropensity считает долю дохода, уходящего на выплату обязательств
// Формула: сумма денег, выплаченных в качестве долгов, за месяц/суммарный располагаемый доход человека за месяц
// Формула преобразования: min{80(0.6-propensity_for_loans); 40}
func (s *Service) LoansPropensity(userID string) (float64, error) {
	q := `
	WITH loans AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total,
		       COALESCE(COUNT(amount_in_rubles), 0) AS count
		FROM expense_in_rubles
		WHERE
			user_id = $1  AND
			type = $2 AND 
			planned = false AND
			date >= NOW() - INTERVAL '30 days'
	),
	income AS (
		SELECT COALESCE(SUM(amount_in_rubles), 0) AS total
		FROM income_in_rubles
		WHERE
			user_id = $1 AND
			planned = false AND
			date >= NOW() - INTERVAL '30 days'
	)
	SELECT 
	    CASE
	        WHEN income.total = 0 THEN 0
			ELSE loans.total / income.total
	    END AS propensity, 
	    loans.count AS count
	FROM loans, income;
	`

	var propensity float64
	var count int
	err := s.repo.QueryRow(q, userID, loan).Scan(&propensity, &count)
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 100, nil
	}

	result := math.Min(80*(0.6-propensity), 40)
	return result, nil
}
