package v1

import (
	"encoding/json"
	"fmt"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
)

type DeltaResponse struct {
	Message    string  `json:"message"`
	Delta      float64 `json:"expenditure_delta"`
	StatusCode int     `json:"status_code"`
}

// ExpenditureDeltaHandler calculates the expenditure delta for an authenticated user.
//
// @Summary Calculate expenditure delta
// @Description This endpoint allows authenticated users to calculate the expenditure delta, providing insight into their financial health.
// @Tags Financial Health
// @Accept  json
// @Produce  json
// @Success 200 {object} DeltaResponse "Successfully calculated expenditure delta"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error while calculating expenditure delta"
// @Security JWT
// @Router /fin_health/expense/delta [get]
func (h *MyHandler) ExpenditureDeltaHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Getting expenditure delta...")

	user, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("auth err"), http.StatusUnauthorized)
		return
	}
	result, err := h.s.FinHealth.ExpenditureDelta(user)
	if err != nil {
		h.errResp(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := DeltaResponse{
		Message:    "Expenditure delta calculated successfully",
		Delta:      result,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

type PropensityResponse struct {
	Message    string  `json:"message"`
	Propensity float64 `json:"expense_propensity"`
	StatusCode int     `json:"status_code"`
}

// ExpensePropensity calculates the expense propensity for an authenticated user.
//
// @Summary Calculate expense propensity
// @Description This endpoint allows authenticated users to calculate the expense propensity, providing insight into their financial health.
// @Tags Financial Health
// @Accept  json
// @Produce  json
// @Success 200 {object} PropensityResponse "Successfully calculated expense propensity"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error while calculating expenditure delta"
// @Security JWT
// @Router /fin_health/expense/propensity [get]
func (h *MyHandler) ExpensePropensity(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Getting expense propensity...")

	user, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("auth err"), http.StatusUnauthorized)
		return
	}
	result, err := h.s.FinHealth.ExpensePropensity(user)
	if err != nil {
		h.errResp(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := PropensityResponse{
		Message:    "expense propensity calculated successfully",
		Propensity: result,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
