package openbanking

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/models"
)

type ConsentRequest struct {
	Permissions        []string `json:"permissions"`
	ExpirationDateTime string   `json:"expirationDateTime"`
}

func GetConsentRequest(auth Auth, apiURL string) (*ConsentRequest, error) {
	url := apiURL + "/account-consents/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+auth.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get consent: status code " + resp.Status)
	}

	var consentRequest ConsentRequest
	if err := json.NewDecoder(resp.Body).Decode(&consentRequest); err != nil {
		return nil, err
	}

	return &consentRequest, nil
}

func CreateConsent(auth Auth, apiURL string, consentReq ConsentRequest) (*models.ConsentResponse, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"Data": consentReq,
	})
	req, err := http.NewRequest("POST", apiURL+"/account-consents", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+auth.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var consentResp models.ConsentResponse
	json.NewDecoder(resp.Body).Decode(&consentResp)
	return &consentResp, nil
}

func GetConsent(auth Auth, consentId, apiURL string) (*models.ConsentResponse, error) {
	// Логика для получения согласия
	return nil, nil
}

func DeleteConsent(auth Auth, consentId, apiURL string) error {
	// Логика для удаления согласия
	return nil
}
