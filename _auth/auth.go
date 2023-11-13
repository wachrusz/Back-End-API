package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Ваш секретный ключ для авторизации по API
const secretAPIKey = "your-secret-key"

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/auth/login", Login).Methods("POST")
	router.HandleFunc("/auth/logout", Logout).Methods("POST")
}

func Login(w http.ResponseWriter, r *http.Request) {
	//TODO: Реализация логики входа
	//TODO: Добавить проверку учетных данных и выдачу токена
	// В данном примере просто возвращаем успешный ответ

	// Добавляем поддержку API-ключа
	apiKey := r.Header.Get("API-Key")
	if apiKey != secretAPIKey {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid API key"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	//TODO: Реализация логики выхода
	//TODO: Добавить логику завершения сеанса пользователя
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}
