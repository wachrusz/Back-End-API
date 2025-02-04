// Package repository provides basic financial repository functionality.
package repository

import (
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
)

type GoalModel struct {
	DB *mydb.Database
}

type GoalTransactionModel struct {
	DB *mydb.Database
}

func (m *GoalModel) Create(goal *models.Goal) (int64, error) {
	var goalID int64
	err := m.DB.QueryRow("INSERT INTO goals (amount, currency_code, user_id, name, months) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		goal.Amount, goal.Currency, goal.UserID, goal.Name, goal.Months).Scan(&goalID)
	if err != nil {
		return 0, err
	}
	return goalID, nil
}

func (m *GoalModel) Update(goal *models.Goal) error {
	result, err := m.DB.Exec(`
		UPDATE goals 
		SET 
			amount = $1,
			currency_code = $2,
            name = $3,
            months = $4,
            additional_months = $5,
            is_completed = $6
		WHERE id = $7
	`, goal.Amount, goal.Currency, goal.Name, goal.Months, goal.AdditionalMonths, goal.IsCompleted, goal.ID)

	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: no goal found with id %d", myerrors.ErrNotFound, goal.ID)
	}

	return nil
}

func (m *GoalModel) Delete(id int64, userID int64) error {
	result, err := m.DB.Exec("DELETE FROM goals WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		// Возвращаем обернутую ошибку, если запрос завершился с ошибкой
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Возвращаем обернутую ошибку при невозможности получить количество строк
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	if rowsAffected == 0 {
		// Возвращаем ошибку, если запись не найдена или не принадлежит пользователю
		return fmt.Errorf("%w: no goal found with id %d for user %d", myerrors.ErrNotFound, id, userID)
	}

	return nil
}

func (m *GoalModel) ListByUserID(userID int64) ([]models.Goal, error) {
	rows, err := m.DB.Query("SELECT id, amount, currency_code, name, months, additional_months, is_completed, start_date FROM goals WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var goal models.Goal
		if err := rows.Scan(&goal.ID, &goal.Amount, &goal.Currency, &goal.Name,
			&goal.Months, &goal.AdditionalMonths, &goal.IsCompleted, &goal.Date); err != nil {
			return nil, err
		}
		goal.UserID = userID
		goals = append(goals, goal)
	}

	return goals, nil
}

func (m *GoalModel) Details(id int64) (*models.GoalDetails, error) {
	var d models.GoalDetails
	q := `
	SELECT
	    -- Цель
		g.amount, 
		g.currency_code, 
		g.user_id, 
		g.name, 
		g.months, 
		g.additional_months, 
		g.is_completed, 
		g.start_date,
		-- Количество месяцев, прошедших с start_date
		EXTRACT(YEAR FROM AGE(CURRENT_DATE, g.start_date)) * 12 + 
		EXTRACT(MONTH FROM AGE(CURRENT_DATE, g.start_date)) AS months_passed,
		-- Общая конвертированная сумма всех транзакций
		COALESCE(SUM(
			CASE
				WHEN gt.currency_code = g.currency_code THEN gt.amount
				ELSE gt.amount * COALESCE(
					(SELECT er.rate_to_ruble
					 FROM exchange_rates er
					 WHERE er.currency_code = gt.currency_code),
					1) / COALESCE(
					(SELECT er.rate_to_ruble
					 FROM exchange_rates er
					 WHERE er.currency_code = g.currency_code),
					1)
			END), 0) AS converted_amount,
		-- Конвертированная сумма транзакций за последний месяц
		COALESCE(SUM(
			CASE
				WHEN gt.created_at >= DATE_TRUNC('month', CURRENT_DATE) AND 
				     gt.created_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month' 
				THEN 
					CASE
						WHEN gt.currency_code = g.currency_code THEN gt.amount
						ELSE gt.amount * COALESCE(
							(SELECT er.rate_to_ruble
							 FROM exchange_rates er
							 WHERE er.currency_code = gt.currency_code),
							1) / COALESCE(
							(SELECT er.rate_to_ruble
							 FROM exchange_rates er
							 WHERE er.currency_code = g.currency_code),
							1)
					END
				ELSE 0
			END), 0) AS last_month_converted_amount
	FROM goals g
	LEFT JOIN goal_transactions gt ON g.id = gt.goal_id AND gt.planned = false
	WHERE g.id = $1
	GROUP BY g.id;
	`

	err := m.DB.QueryRow(q, id).Scan(d.Goal.Amount, d.Goal.Currency, d.Goal.UserID, d.Goal.Name,
		d.Goal.Months, d.Goal.AdditionalMonths, d.Goal.IsCompleted, d.Goal.Date,
		&d.Month, &d.Gathered, &d.CurrentPayment)
	if err != nil {
		return nil, err
	}

	d.Goal.ID = id
	d.MonthlyPayment = d.Goal.Amount / float64(d.Goal.Months)
	d.CurrentNeed = d.MonthlyPayment - d.CurrentPayment

	return &d, nil
}

func (m *GoalTransactionModel) Create(transaction *models.GoalTransaction) (int64, error) {
	var transactionId int64
	err := m.DB.QueryRow("INSERT INTO goal_transactions(goal_id, amount, planned, currency_code, connected_account) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		transaction.GoalID, transaction.Amount, transaction.Planned, transaction.Currency, transaction.BankAccount).Scan(&transactionId)

	if err != nil {
		return 0, err
	}

	return transactionId, nil
}

func (m *GoalTransactionModel) ListByGoal(goalID int64) ([]*models.GoalTransaction, error) {
	rows, err := m.DB.Query("SELECT id, amount, planned, currency_code, date, connected_account FROM goal WHERE goal_id = $1", goalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.GoalTransaction
	for rows.Next() {
		var transaction models.GoalTransaction
		transaction.GoalID = goalID

		if err := rows.Scan(&transaction.ID, &transaction.Amount, &transaction.Planned, &transaction.Currency, &transaction.Date, &transaction.BankAccount); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}
