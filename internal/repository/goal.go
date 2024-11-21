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

func (m *GoalModel) Create(goal *models.Goal) (int64, error) {
	var goalID int64
	err := m.DB.QueryRow("INSERT INTO goal (goal, need, current_state, end_date, currency, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		goal.Goal, goal.Need, goal.CurrentState, goal.EndDate, goal.Currency, goal.UserID).Scan(&goalID)
	if err != nil {
		return 0, err
	}
	return goalID, nil
}

func (m *GoalModel) Update(goal *models.Goal) error {
	result, err := m.DB.Exec(`
		UPDATE goal 
		SET 
			goal = $1,
			need = $2,
			end_date = $3,
			currency = $4
		WHERE id = $5 and user_id = $6
	`, goal.Goal, goal.Need, goal.EndDate, goal.Currency, goal.ID, goal.UserID)

	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: no goal found with id %s for user %s", myerrors.ErrNotFound, goal.ID, goal.UserID)
	}

	return nil
}

func (m *GoalModel) Delete(id, userID string) error {
	result, err := m.DB.Exec("DELETE FROM goal WHERE id = $1 AND user_id = $2", id, userID)
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
		return fmt.Errorf("%w: no goal found with id %s for user %s", myerrors.ErrNotFound, id, userID)
	}

	return nil
}

func (m *GoalModel) GetByUserID(userID string) ([]models.Goal, error) {
	rows, err := m.DB.Query("SELECT id, goal, need, current_state, start_date, end_date, currency FROM goal WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var goal models.Goal
		if err := rows.Scan(&goal.ID, &goal.Goal, &goal.Need, &goal.CurrentState, goal.StartDate, &goal.EndDate, &goal.Currency); err != nil {
			return nil, err
		}
		goal.UserID = userID
		goals = append(goals, goal)
	}

	return goals, nil
}
