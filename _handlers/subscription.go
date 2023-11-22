package handlers

import (
	models "backEndAPI/_models"
	"encoding/json"
	"net/http"
)

// CreateSubscriptionHandler обрабатывает запрос на создание категории расходов.
func CreateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	var subscription models.Subscription

	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := models.CreateSubscription(&subscription)
	if err != nil {
		http.Error(w, "Error creating subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Subscription created successfully"))
}
