//go:build !exclude_swagger
// +build !exclude_swagger

// Package handlers provides http functionality.
package v1

import (
	"encoding/json"
	"errors"
	"github.com/wachrusz/Back-End-API/internal/auth"
	"github.com/wachrusz/Back-End-API/internal/models"
	"net/http"
)

// @Summary Create an income
// @Description Create a new income.
// @Tags Analytics
// @Accept json
// @Produce json
// @Param income body models.Income true "Income object"
// @Success 201 {string} string "Income created successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error creating income"
// @Security JWT
// @Router /analytics/income [post]
func CreateIncomeHandler(w http.ResponseWriter, r *http.Request) {
	var income models.Income
	if err := json.NewDecoder(r.Body).Decode(&income); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	income.UserID = userID

	incomeID, err := models.CreateIncome(&income)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error creating income: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":           "Successfully created an income",
		"created_object_id": incomeID,
		"status_code":       http.StatusCreated,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}
