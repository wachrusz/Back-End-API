//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package auth

import (
	logger "backEndAPI/_logger"
	mydb "backEndAPI/_mydatabase"
	"net"

	"fmt"
	"net/http"
	"sync"
)

type ActiveUser struct {
	UserID   string
	Username string
	DeviceID string
}

var (
	activeUsers = make(map[string]ActiveUser)
	activeMu    sync.Mutex
)

// @Title Get active user sessions
// @Description Get information about active user sessions.
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]ActiveUser "List of active user sessions"
// @Router /auth/sessions [get]
func InitActiveUsers() {
	var query string = "SELECT user_id, username, device_id FROM sessions"
	rows, err := mydb.GlobalDB.Query(query)
	if err != nil {
		logger.ErrorLogger.Printf("Unnable to check DB DUE TO: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userID, username, deviceID string
		if err := rows.Scan(&userID, &username, &deviceID); err != nil {
			logger.ErrorLogger.Printf("Strange error at: %v", err)
		}

		activeUser := ActiveUser{
			UserID:   userID,
			Username: username,
			DeviceID: deviceID,
		}

		activeUsers[userID] = activeUser
	}
	for userID, activeUser := range activeUsers {
		fmt.Printf("User ID: %s, Username: %s, Device ID: %s\n", userID, activeUser.Username, activeUser.DeviceID)
	}
}

func getActiveUser(userID string) ActiveUser {
	return activeUsers[userID]
}

func IsUserActive(userID string) bool {
	activeMu.Lock()
	defer activeMu.Unlock()

	_, ok := activeUsers[userID]
	return ok
}

func AddActiveUser(userID, username, deviceID string) {
	activeMu.Lock()
	defer activeMu.Unlock()

	activeUsers[userID] = ActiveUser{userID, username, deviceID}
}

func RemoveActiveUser(userID string) {
	activeMu.Lock()
	defer activeMu.Unlock()

	delete(activeUsers, userID)
}

/*
func SetDeviceIDInSession(r *http.Request, w http.ResponseWriter, deviceID string) {
	defer sessionMutex.Unlock()
	sessionMutex.Lock()

	//! TODO СЕССИОНИРОВАНИЕ СРОЧНО
	if err != nil {
		logger.ErrorLogger.Printf("Error saving session: %v", err)
		return
	}
}*/

// DATABASE OPERATIONS
func saveSessionToDatabase(username, deviceID, user_id string) error {
	_, err := mydb.GlobalDB.Exec(`
        INSERT INTO sessions (username, device_id, created_at, last_activity, user_id)
        VALUES ($1, $2, NOW(), NOW(), $3)`,
		username, deviceID, user_id)
	return err
}

func checkSessionInDatabase(username, deviceID string) (bool, error) {
	var count int

	err := mydb.GlobalDB.QueryRow(`
        SELECT COUNT(*) FROM sessions WHERE
            username = $1 AND device_id = $2;`,
		username, deviceID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking session in database: %v", err)
	}

	return count > 0, nil
}

func removeSessionFromDatabase(deviceID string) error {
	_, err := mydb.GlobalDB.Exec(`
        DELETE FROM sessions WHERE device_id = $1`,
		deviceID)
	return err
}

// SERVICE FUNCTIONS
func isDeviceIDAlreadyUsed(db *mydb.Database, username, deviceID string) (error, bool) {

	query := "SELECT COUNT(*) FROM sessions WHERE username = $1 AND device_id = $2"

	var count int
	err := db.QueryRow(query, username, deviceID).Scan(&count)
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
	SELECT id FROM users WHERE username = $1;
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
