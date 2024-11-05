package openbanking

import (
	"bytes"
	"errors"
	"net/http"
)

type OpenBankingClient struct {
	Auth    Auth
	BaseURL string
}

func NewClient(clientID, clientSecret, baseURL string) (*OpenBankingClient, error) {
	auth := Auth{}
	err := auth.GetToken(clientID, clientSecret, baseURL+"/token")
	if err != nil {
		return nil, err
	}

	return &OpenBankingClient{
		Auth:    auth,
		BaseURL: baseURL,
	}, nil
}

func (client *OpenBankingClient) MakeRequest(method, endpoint string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, client.BaseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+client.Auth.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New("API request failed with status: " + resp.Status)
	}

	return resp, nil
}
