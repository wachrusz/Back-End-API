package fin_health

import "math"

func (s *Service) ExpenditureDelta(userID string) (float64, error) {
	q := `
	WITH current_month_expenses AS (
		SELECT COALESCE(SUM(
			CASE
				WHEN currency_code = 'RUB' THEN amount
				ELSE amount * COALESCE(
					(SELECT rate_to_ruble 
					 FROM exchange_rates 
					 WHERE currency_code = expense.currency_code), 
					1)
			END), 0) AS total_current_month
		FROM expense
		WHERE
			user_id = $1 AND
			date >= NOW() - INTERVAL '30 days'
	),
	average_monthly_expenses AS (
		SELECT COALESCE(SUM(
			CASE
				WHEN currency_code = 'RUB' THEN amount
				ELSE amount * COALESCE(
					(SELECT rate_to_ruble 
					 FROM exchange_rates 
					 WHERE currency_code = expense.currency_code), 
					1)
			END), 0) / 3 AS avg_last_3_months
		FROM expense
		WHERE
			user_id = $1 AND
			date >= NOW() - INTERVAL '90 days'
	)
	SELECT 
		(cm.total_current_month - am.avg_last_3_months) / am.avg_last_3_months * 100 AS delta
	FROM 
		current_month_expenses cm, 
		average_monthly_expenses am;
	`

	var delta float64
	err := s.repo.QueryRow(q, userID).Scan(&delta)
	if err != nil {
		return 0, err
	}

	result := math.Min(2.5*(15.0-delta), 50)
	return result, nil
}