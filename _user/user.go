package user

import (
	"database/sql"
	"errors"
	"sync"

	logger "backEndAPI/_logger"
	mydb "backEndAPI/_mydatabase"

	"golang.org/x/crypto/bcrypt"
)

// User - структура для представления пользователя
type User struct {
	Username       string
	HashedPassword string
}

// mutex - мьютекс для безопасного доступа к мапе Users
var (
	mutex sync.Mutex
)

// RegisterUser - функция для регистрации нового пользователя
func RegisterUser(username, name, password string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := GetUserByUsername(username); exists {
		errMsg := "User with username " + username + " already exists"
		logger.ErrorLogger.Println(errMsg)
		return errors.New("Already exists")
	}

	if username == "" || name == "" || password == "" {
		return errors.New("Blank fields are not allowed")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = mydb.GlobalDB.Exec("INSERT INTO users (username, hashed_password, name) VALUES ($1, $2, $3)", username, hashedPassword, name)
	if err != nil {
		logger.ErrorLogger.Println("Error inserting user:", err)
		return err
	}

	logger.InfoLogger.Printf("New user registered: %s\n", username)

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

func GetUserByUsername(username string) (User, bool) {
	var user User
	var id int

	row := mydb.GlobalDB.QueryRow("SELECT id, username, hashed_password FROM users WHERE username = $1", username)
	err := row.Scan(&id, &user.Username, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return user, false
	} else if err != nil {
		logger.ErrorLogger.Println("Error querying user:", err)
		return user, false
	}
	return user, true
}

func GetHashedPasswordByUsername(username string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	user, exists := GetUserByUsername(username)
	if !exists {
		errMsg := "User with username " + username + " not found"
		logger.ErrorLogger.Println(errMsg)
		return "", errors.New("User not found")
	}

	return user.HashedPassword, nil
}
