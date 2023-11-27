//go:build !exclude_swagger
// +build !exclude_swagger

// Package history provides operations archive functionality.
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

// @Summary Get history entries
// @Description Get a list of history entries.
// @Tags History
// @Produce json
// @Success 200 {array} string "List of history entries"
// @Router /history/get [get]
func GetHistory(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики получения истории

	history := []string{}
	json.NewEncoder(w).Encode(history)
}

// @Summary Add history entry
// @Description Add a new entry to the history.
// @Tags History
// @Accept json
// @Produce json
// @Success 200 {string} string "History entry added successfully"
// @Router /history/add [post]
func AddHistoryEntry(w http.ResponseWriter, r *http.Request) {
	//TODO Реализация логики добавления записи в историю

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("History entry added successfully"))
}
