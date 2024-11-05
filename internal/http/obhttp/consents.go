package obhttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/openbanking"
)

// CreateConsentHandler handles requests to create consent for accessing user accounts from the Open Banking API.
//
// @Summary Create consent
// @Description Initiates the consent process for accessing user accounts using the provided authorization token and API URL.
// @Tags OpenBanking
// @Accept json
// @Produce json
// @Param Auth body openbanking.Auth true "Authorization token for Open Banking"
// @Param api_url query string true "API base URL for Open Banking"
// @Success 200 {object} map[string]interface{} "Successfully created consent"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload or missing URL parameter"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal service error while creating consent"
// @Router /openbanking/consent/create [post]
func (h *MyHandler) CreateConsentHandler(w http.ResponseWriter, r *http.Request) {
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

	consentRequest, err := openbanking.GetConsentRequest(account, url)
	if err != nil {
		switch err {
		default:
			h.errResp(w, fmt.Errorf("Internal service error"), http.StatusInternalServerError)
			return
		}
	}

	openbanking.CreateConsent(account, url, *consentRequest)
}

// GetConsentHandler handles requests to fetch consent information from the Open Banking API.
//
// @Summary Retrieve consent
// @Description Fetches consent information using the provided authorization token, consent ID, and API URL.
// @Tags OpenBanking
// @Accept json
// @Produce json
// @Param Auth body openbanking.Auth true "Authorization token for Open Banking"
// @Param api_url query string true "API base URL for Open Banking"
// @Param consent_id query string true "Consent ID to retrieve"
// @Success 200 {object} map[string]interface{} "Successfully fetched consent"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload, missing parameters, or incorrect ID"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal service error while fetching consent"
// @Router /openbanking/consent/get [get]
func (h *MyHandler) GetConsentHandler(w http.ResponseWriter, r *http.Request) {
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
	consentID := r.URL.Query().Get("consent_id")
	if url == "" {
		h.errResp(w, fmt.Errorf("Incorrect id"), http.StatusBadRequest)
		return
	}

	consent, err := openbanking.GetConsent(account, consentID, url)
	if err != nil {
		switch err {
		default:
			h.errResp(w, fmt.Errorf("Internal service error"), http.StatusInternalServerError)
			return
		}
	}

	//TEMP
	response := map[string]interface{}{
		"message":     "Successfully fetched consent.",
		"consent":     consent,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

// DeleteConsentHandler handles requests to delete consent from the Open Banking API.
//
// @Summary Delete consent
// @Description Deletes existing consent for accessing user accounts using the provided authorization token, consent ID, and API URL.
// @Tags OpenBanking
// @Accept json
// @Produce json
// @Param Auth body openbanking.Auth true "Authorization token for Open Banking"
// @Param api_url query string true "API base URL for Open Banking"
// @Param consent_id query string true "Consent ID to delete"
// @Success 200 {object} map[string]interface{} "Successfully deleted consent"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload, missing parameters, or incorrect ID"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal service error while deleting consent"
// @Router /openbanking/consent/delete [delete]
func (h *MyHandler) DeleteConsentHandler(w http.ResponseWriter, r *http.Request) {
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
	consentID := r.URL.Query().Get("consent_id")
	if url == "" {
		h.errResp(w, fmt.Errorf("Incorrect id"), http.StatusBadRequest)
		return
	}

	err := openbanking.DeleteConsent(account, consentID, url)
	if err != nil {
		switch err {
		default:
			h.errResp(w, fmt.Errorf("Internal service error"), http.StatusInternalServerError)
			return
		}
	}

	//TEMP
	response := map[string]interface{}{
		"message":     "Successfully deleted consent.",
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}
