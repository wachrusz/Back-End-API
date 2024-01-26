//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	enc "backEndAPI/_encryption"
	logger "backEndAPI/_logger"
	mydb "backEndAPI/_mydatabase"

	"net"

	"fmt"
	"net/http"
	"sync"
)

type ActiveUser struct {
	UserID         string
	Email          string
	DeviceID       string
	DecryptedToken string
}

var (
	ActiveUsers = make(map[string]ActiveUser)
	activeMu    sync.Mutex
)

func InitActiveUsers() {
	var query string = "SELECT user_id, email, device_id, token FROM sessions"
	rows, err := mydb.GlobalDB.Query(query)
	if err != nil {
		logger.ErrorLogger.Printf("Unnable to check DB DUE TO: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID, email, deviceID, token string
		if err := rows.Scan(&userID, &email, &deviceID, &token); err != nil {
			logger.ErrorLogger.Printf("Strange error at: %v", err)
		}

		decryptedToken, err := enc.DecryptToken(token)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to decrypt token for UserID: %v", userID, ", token: %v", decryptedToken)
		}

		activeUser := ActiveUser{
			UserID:         userID,
			Email:          email,
			DeviceID:       deviceID,
			DecryptedToken: decryptedToken,
		}

		ActiveUsers[userID] = activeUser
	}
}

func getActiveUser(userID string) ActiveUser {
	return ActiveUsers[userID]
}

func IsUserActive(userID string) bool {
	activeMu.Lock()
	defer activeMu.Unlock()

	_, ok := ActiveUsers[userID]
	if !ok {
		logger.ErrorLogger.Printf("User %s is not active", userID)
	}
	return ok
}

func AddActiveUser(userID, email, deviceID, token string) {
	activeMu.Lock()
	defer activeMu.Unlock()

	ActiveUsers[userID] = ActiveUser{userID, email, deviceID, token}
}

func RemoveActiveUser(userID string) {
	activeMu.Lock()
	defer activeMu.Unlock()

	delete(ActiveUsers, userID)
}

// *NEW
func SetAccessToken(userID, newAccessToken string) {
	activeMu.Lock()
	defer activeMu.Unlock()

	if ActiveUsers == nil {
		ActiveUsers = make(map[string]ActiveUser)
	}

	user, exists := ActiveUsers[userID]
	if !exists {
		user = ActiveUser{}
		ActiveUsers[userID] = user
	}

	user.DecryptedToken = newAccessToken
}

// DATABASE OPERATIONS
func saveSessionToDatabase(email, deviceID, user_id, token string) error {

	encryptedToken, err := enc.EncryptToken(token)
	if err != nil {
		return err
	}

	_, err = mydb.GlobalDB.Exec(`
        INSERT INTO sessions (email, device_id, created_at, last_activity, user_id, token)
        VALUES ($1, $2, NOW(), NOW(), $3, $4)`,
		email, deviceID, user_id, encryptedToken)
	return err
}

// *NEW
func updateLastActivity(userID string) error {
	query := `
	UPDATE sessions
	SET last_activity = NOW()
	WHERE user_id = $1;
	`

	_, err := mydb.GlobalDB.Exec(query, userID)
	return err
}

func checkSessionInDatabase(email, deviceID string) (bool, error) {
	var count int

	err := mydb.GlobalDB.QueryRow(`
        SELECT COUNT(*) FROM sessions WHERE
            email = $1 AND device_id = $2;`,
		email, deviceID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking session in database: %v", err)
	}

	return count > 0, nil
}

func removeSessionFromDatabase(deviceID, userID string) error {
	_, err := mydb.GlobalDB.Exec(`
        DELETE FROM sessions WHERE device_id = $1 AND user_id = $2`,
		deviceID, userID)
	return err
}

// SERVICE FUNCTIONS
func isDeviceIDAlreadyUsed(db *mydb.Database, email, deviceID string) (error, bool) {

	query := "SELECT COUNT(*) FROM sessions WHERE email = $1 AND device_id = $2"

	var count int
	err := db.QueryRow(query, email, deviceID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking device ID: %v", err), false
	}

	return nil, count > 0
}

// !FIX
func GetDeviceIDFromRequest(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	deviceID := fmt.Sprintf("%s_%s", ip, r.UserAgent())

	return deviceID
}

func getUserIDFromUsersDatabase(usernameOrDeviceID string) (string, error) {
	var result string

	err := mydb.GlobalDB.QueryRow(`
	SELECT id FROM users WHERE email = $1;
	`, usernameOrDeviceID).Scan(&result)

	if err != nil {
		return "", fmt.Errorf("error checking session in database: %v", err)
	}
	return result, nil

}

func GetUserIDFromSessionDatabase(usernameOrDeviceID string) (string, error) {
	var result string
	err := mydb.GlobalDB.QueryRow(`
	SELECT user_id FROM sessions WHERE device_id = $1;
	`, usernameOrDeviceID).Scan(&result)
	if err != nil {
		return "", fmt.Errorf("error checking session in database: %v", err)
	}
	return result, nil
}
