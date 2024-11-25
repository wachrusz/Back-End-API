package openbanking

import (
	"encoding/json"
	"net/http"

	"github.com/wachrusz/Back-End-API/internal/repository"
)

type AccountListResponse struct {
	Accounts []repository.Account `json:"accounts"`
}

func GetAccounts(auth Auth, apiURL string) (*AccountListResponse, error) {
	req, err := http.NewRequest("GET", apiURL+"/accounts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+auth.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accountListResp AccountListResponse
	json.NewDecoder(resp.Body).Decode(&accountListResp)
	return &accountListResp, nil
}
