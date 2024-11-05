package openbanking

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Auth struct {
	Token string
}

func (a *Auth) GetToken(clientID, clientSecret, authURL string) error {
	reqBody := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", authURL, reqBody)
	if err != nil {
		return err
	}

	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	a.Token = result["access_token"].(string)
	return nil
}
