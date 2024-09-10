//go:build !exclude_swagger
// +build !exclude_swagger

// Package user provides register functionality.
package user

import (
	"database/sql"
	"errors"
	"main/pkg/logger"
	mydb "main/pkg/mydatabase"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email          string
	HashedPassword string
}

var (
	mutex sync.Mutex
)

func RegisterUser(email, password string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := GetUserByEmail(email); exists {
		errMsg := "User with email " + email + " already exists"
		logger.ErrorLogger.Println(errMsg)
		return errors.New("Already exists")
	}

	if email == "" || password == "" {
		return errors.New("Blank fields are not allowed")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = mydb.GlobalDB.Exec("INSERT INTO users (email, hashed_password) VALUES ($1, $2)", email, hashedPassword)
	if err != nil {
		logger.ErrorLogger.Println("Error inserting user:", err)
		return err
	}

	logger.InfoLogger.Printf("New user registered: %s\n", email)

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLogger.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func GetUserByEmail(email string) (User, bool) {
	var user User
	var id int

	row := mydb.GlobalDB.QueryRow("SELECT id, email, hashed_password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &user.Email, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return user, false
	} else if err != nil {
		logger.ErrorLogger.Println("Error querying user:", err)
		return user, false
	}
	return user, true
}

func GetHashedPasswordByUsername(email string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	user, exists := GetUserByEmail(email)
	if !exists {
		errMsg := "User with email " + email + " not found"
		logger.ErrorLogger.Println(errMsg)
		return "", errors.New("User not found")
	}

	return user.HashedPassword, nil
}
