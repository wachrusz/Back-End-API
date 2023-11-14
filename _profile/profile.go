package profile

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Простой тип для представления профиля пользователя
type UserProfile struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

// userProfile - имитация базы данных
var userProfile = make(map[string]UserProfile)

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/profile/get", GetProfile).Methods("GET")
	router.HandleFunc("/profile/update", UpdateProfile).Methods("PUT")
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	// Реализация логики получения профиля
	username := "user1" // В данном примере захардкодим пользователя
	profile, ok := userProfile[username]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Profile not found"))
		return
	}

	json.NewEncoder(w).Encode(profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики обновления профиля

	// В данном примере просто возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Profile updated successfully"))
}
