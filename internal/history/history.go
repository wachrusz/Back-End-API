//go:build !exclude_swagger
// +build !exclude_swagger

// Package history provides operations archive functionality.
package history

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func RegisterHandlers(router chi.Router) {
	router.Get("/history/get", GetHistory)
	router.Post("/history/add", AddHistoryEntry)
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

// @Summary Create history entry
// @Description Create a new entry to the history.
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
