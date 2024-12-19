package v1

import (
	"encoding/json"
	"errors"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
)

type DeltaResponse struct {
	Message    string  `json:"message"`
	Delta      float64 `json:"expenditure_delta"`
	StatusCode int     `json:"status"`
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
		h.errResp(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}
	result, err := h.s.FinHealth.ExpenditureDelta(user)
	if err != nil {
		h.errResp(w, err, http.StatusInternalServerError)
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
