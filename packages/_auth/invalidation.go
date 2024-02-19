package auth

import (
	mydb "main/packages/_mydatabase"
)

func InvalidateTokensByUserID(userID string) error {
	_, err := mydb.GlobalDB.Exec(`DELETE FROM sessions WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	return nil
}

/*
Хорошим вариантом может быть реализация инвалидации посредством добавления атрибута просроченности токену,
либо добавление бан листа в БД, хз как лучше. Пока что лучше оставить просто удаление, МВП же)))
*/
