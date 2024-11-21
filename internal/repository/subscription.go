// Package repository provides basic financial repository functionality.
package repository

import (
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	"time"
)

type SubscriptionModel struct {
	DB *mydb.Database
}

func (m *SubscriptionModel) Create(subscription *models.Subscription) (int64, error) {
	parsedDateStart, err := time.Parse("2006-01-02", subscription.StartDate)
	if err != nil {
		return 0, err
	}

	parsedDateEnd, err1 := time.Parse("2006-01-02", subscription.EndDate)
	if err1 != nil {
		return 0, err1
	}

	var subscriptionID int64
	err2 := m.DB.QueryRow("INSERT INTO subscriptions (user_id, start_date, end_date, is_active) VALUES ($1, $2, $3, $4) RETURNING id",
		subscription.UserID, parsedDateStart, parsedDateEnd, subscription.IsActive).Scan(&subscriptionID)
	if err2 != nil {
		return 0, err2
	}
	return subscriptionID, nil
}

func (m *SubscriptionModel) Update(subscription *models.Subscription) error {
	parsedDateStart, err := time.Parse("2006-01-02", subscription.StartDate)
	if err != nil {
		return err
	}

	parsedDateEnd, err1 := time.Parse("2006-01-02", subscription.EndDate)
	if err1 != nil {
		return err1
	}

	result, err := m.DB.Exec(`
		UPDATE subscriptions SET 
			start_date = $1,
			end_date = $2,
			is_active = $3
		WHERE id = $4 and user_id = $5`,
		parsedDateStart, parsedDateEnd, subscription.IsActive, subscription.ID, subscription.UserID)

	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	// Проверяем, сколько строк было затронуто
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrInternal, err) // Ошибка получения числа затронутых строк
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: no subscription found with id %s for user %s", myerrors.ErrNotFound, subscription.ID, subscription.UserID)
	}

	return nil
}

func (m *SubscriptionModel) Delete(id, userID string) error {
	result, err := m.DB.Exec("DELETE FROM subscriptions WHERE id = $1 AND user_id = $2", id, userID)
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
		return fmt.Errorf("%w: no subscription found with id %s for user %s", myerrors.ErrNotFound, id, userID)
	}

	return nil
}
