package categories

import (
	//"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/categories/analytics", GetAnalytics).Methods("GET")
	router.HandleFunc("/categories/tracker", GetTracker).Methods("GET")
	router.HandleFunc("/categories/settings", GetSettings).Methods("GET")
}

func GetAnalytics(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики получения данных по категории "Аналитика"

	// В данном примере просто возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Analytics data"))
}

func GetTracker(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики получения данных по категории "Трекер"

	// В данном примере просто возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tracker data"))
}

func GetSettings(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики получения данных по категории "Настройки"
	// В данном примере просто возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Settings data"))
}
