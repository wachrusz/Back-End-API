//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package user

import (
	"database/sql"
	"fmt"

	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
)

type ActiveUser struct {
	UserID         string
	Email          string
	DeviceID       string
	DecryptedToken string
}

func (s *Service) InitActiveUsers() {
	var query = "SELECT user_id, email, device_id, token FROM sessions"
	rows, err := s.repo.Query(query)
	if err != nil {
	}
	defer rows.Close()

	for rows.Next() {
		var userID, email, deviceID, token string
		if err := rows.Scan(&userID, &email, &deviceID, &token); err != nil {
		}

		decryptedToken, err := enc.DecryptToken(token)
		_ = decryptedToken // FIXME: хз что тут
		if err != nil {
		}
	}
}

func (s *Service) GetActiveUser(userID string) ActiveUser {
	return s.ActiveUsers[userID]
}

func (s *Service) IsUserActive(userID string) bool {
	s.activeMu.Lock()
	defer s.activeMu.Unlock()

	_, ok := s.ActiveUsers[userID]
	if !ok {

	query := `
		SELECT 1 FROM sessions
		WHERE user_id = $1;
	`
	row := s.repo.QueryRow(query, userID)
	var dummy int
	err := row.Scan(&dummy)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Printf("User %s is not active", userID)
			return false
		}
		logger.ErrorLogger.Printf("Error executing query for user %s: %v", userID, err)
		return false

	}
	return true
}

func (s *Service) AddActiveUser(user ActiveUser) {
	s.activeMu.Lock()
	defer s.activeMu.Unlock()

	s.ActiveUsers[user.UserID] = user
}

func (s *Service) RemoveActiveUser(userID string) {
	s.activeMu.Lock()
	defer s.activeMu.Unlock()

	delete(s.ActiveUsers, userID)
}

// DATABASE OPERATIONS
func (s *Service) SaveSessionToDatabase(email, deviceID, userID, token string) error {

	encryptedToken, err := enc.EncryptToken(token)
	if err != nil {
		return err
	}

	_, err = s.repo.Exec(`
        INSERT INTO sessions (email, device_id, created_at, last_activity, user_id, token)
        VALUES ($1, $2, NOW(), NOW(), $3, $4)`,
		email, deviceID, userID, encryptedToken)
	return err
}

// *NEW
func (s *Service) UpdateLastActivity(userID string) error {
	query := `
	UPDATE sessions
	SET last_activity = NOW()
	WHERE user_id = $1;
	`

	_, err := s.repo.Exec(query, userID)
	return err
}

func (s *Service) CheckSessionInDatabase(email, deviceID string) (bool, error) {
	var count int

	err := s.repo.QueryRow(`
        SELECT COUNT(*) FROM sessions WHERE
            email = $1 AND device_id = $2;`,
		email, deviceID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking session in database: %v", err)
	}

	return count > 0, nil
}

func (s *Service) RemoveSessionFromDatabase(deviceID, userID string) error {
	_, err := s.repo.Exec(`
        DELETE FROM sessions WHERE device_id = $1 AND user_id = $2`,
		deviceID, userID)
	return err
}

// SERVICE FUNCTIONS
func (s *Service) IsDeviceIDAlreadyUsed(email, deviceID string) (error, bool) {

	query := "SELECT COUNT(*) FROM sessions WHERE email = $1 AND device_id = $2"

	var count int
	err := s.repo.QueryRow(query, email, deviceID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking device ID: %v", err), false
	}

	return nil, count > 0
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

func (s *Service) GetUserIDFromSessionDatabase(usernameOrDeviceID string) (string, error) {
	var result string
	err := s.repo.QueryRow(`
	SELECT user_id FROM sessions WHERE device_id = $1;
	`, usernameOrDeviceID).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("error checking session in database: %v", err)
	}
	return result, nil
}

func (s *Service) SetAccessToken(userID, newAccessToken string) {
	s.activeMu.Lock()
	defer s.activeMu.Unlock()

	if s.ActiveUsers == nil {
		s.ActiveUsers = make(map[string]ActiveUser)
	}

	user, exists := s.ActiveUsers[userID]
	if !exists {
		user = ActiveUser{}
		s.ActiveUsers[userID] = user
	}

	user.DecryptedToken = newAccessToken
}
