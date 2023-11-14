package auth

import (
	"net/http"

	logger "backEndAPI/_logger"
	mydb "backEndAPI/_mydatabase"
	user "backEndAPI/_user"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var (
	globalDB     *mydb.Database
	secretAPIKey string
)

func SetDB(db *mydb.Database) {
	globalDB = db
}

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/auth/login", Login).Methods("POST")
	router.HandleFunc("/auth/logout", Logout).Methods("POST")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		username string
		password string
	)

	username = r.FormValue("username")
	password = r.FormValue("password")

	// Проверка наличия логина и пароля в запросе
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing username or password"))
		logger.ErrorLogger.Printf("Missing username or password in login request from %s\n", r.RemoteAddr)
		return
	}

	// Проверка правильности логина и пароля (реализуйте свою логику проверки)
	if !isValidCredentials(username, password) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid username or password"))
		logger.ErrorLogger.Printf("Invalid username or password in login request from %s\n", r.RemoteAddr)
		return
	}

	// Добавляем поддержку API-ключа
	apiKey := r.Header.Get("API-Key")
	if apiKey != secretAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API key"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))

	logger.InfoLogger.Printf("User %s logged in from %s\n", username, r.RemoteAddr)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики выхода
	//TODO Добавить логику завершения сеанса пользователя
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

func SetAPIKey(apiKey string) {
	secretAPIKey = apiKey
}

func GetAPIKey() string {
	return secretAPIKey
}

func isValidCredentials(username, password string) bool {
	hashedPassword, ok := user.GetHashedPasswordByUsername(username)
	if ok != nil {
		return false
	}
	if comparePasswords(hashedPassword, password) {
		return true
	}
	return false
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
