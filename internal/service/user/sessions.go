package user

import (
	"fmt"

	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
)

func (s *Service) SaveSessionToDatabase(userID, deviceID, token string) error {
	encryptedToken, err := enc.EncryptToken(token)
	if err != nil {
		return err
	}

	_, err = s.repo.Exec(`
        INSERT INTO sessions (device_id, created_at, last_activity, user_id, token, expires_at)
        VALUES ($1, NOW(), NOW(), $2, $3, NOW() + INTERVAL '15 minutes')`, deviceID, userID, encryptedToken)
	return err
}

func (s *Service) RemoveSessionFromDatabase(deviceID, userID string) error {
	// TODO: remake
	_, err := s.repo.Exec(`
        DELETE FROM sessions WHERE device_id = $1 AND user_id = $2`,
		deviceID, userID)
	return err
}

func (s *Service) GetUserIDFromUsersDatabase(usernameOrDeviceID string) (string, error) {
	var result string

	err := s.repo.QueryRow(`
	SELECT id FROM users WHERE email = $1;
	`, usernameOrDeviceID).Scan(&result)

	if err != nil {
		return "", fmt.Errorf("error checking session in database: %v", err)
	}
	return result, nil
}
