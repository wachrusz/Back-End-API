package obhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/openbanking"
)

// GetAccountsHandler handles requests to fetch a list of user accounts from the Open Banking API.
//
// @Summary Retrieve user accounts
// @Description Fetches a list of user accounts from the Open Banking API using the provided authorization token and API URL.
// @Tags OpenBanking
// @Accept json
// @Produce json
// @Param Auth body openbanking.Auth true "Authorization token for Open Banking"
// @Param api_url query string true "API base URL for Open Banking"
// @Success 200 {object} map[string]interface{} "Successfully fetched accounts"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload or missing URL parameter"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal service error while fetching accounts"
// @Router /openbanking/accounts/get [get]
func (h *MyHandler) GetAccountsHandler(w http.ResponseWriter, r *http.Request) {
	var account openbanking.Auth
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	url := r.URL.Query().Get("api_url")
	if url == "" {
		h.errResp(w, fmt.Errorf("Incorrect url"), http.StatusBadRequest)
		return
	}

	accounts, err := openbanking.GetAccounts(account, url)
	if err != nil {
		switch err {
		default:
			h.errResp(w, fmt.Errorf("Internal service error"), http.StatusInternalServerError)
			return
		}
	}
	response := map[string]interface{}{
		"message":     "Successfully fetched accounts.",
		"accounts":    accounts,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}
