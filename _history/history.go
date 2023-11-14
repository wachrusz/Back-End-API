package history

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/history/get", GetHistory).Methods("GET")
	router.HandleFunc("/history/add", AddHistoryEntry).Methods("POST")
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики получения истории

	// В данном примере просто возвращаем пустой массив
	history := []string{}
	json.NewEncoder(w).Encode(history)
}

func AddHistoryEntry(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики добавления записи в историю
	// В данном примере просто возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("History entry added successfully"))
}
