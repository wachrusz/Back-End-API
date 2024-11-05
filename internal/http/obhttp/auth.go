package obhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/models"
	"github.com/wachrusz/Back-End-API/internal/openbanking"
)

// GetTokenHandler handles requests to obtain an authorization token from the external banking service.
//
// @Summary Obtain authorization token
// @Description Retrieves an authorization token from the external banking service using client credentials.
// @Tags OpenBanking
// @Accept json
// @Produce json
// @Param GetTokenRequest body models.GetTokenRequest true "Request object containing client credentials"
// @Success 200 {object} openbanking.Auth "Token retrieved successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 500 {object} jsonresponse.Error
func (h *MyHandler) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Getting Token from external bank")

	var request models.GetTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}
	var token openbanking.Auth
	err := token.GetToken(request.ClientID, request.ClientSecret, request.AuthURL)
	if err != nil {
		switch err {
		default:
			h.errResp(w, fmt.Errorf("internal error"), http.StatusInternalServerError)
			return
		}
	}
}
