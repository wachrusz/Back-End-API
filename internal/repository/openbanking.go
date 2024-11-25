package repository

type GetTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AuthURL      string `json:"auth_url"`
}

type Consent struct {
	Permissions             []string `json:"permissions"`
	ExpirationDateTime      string   `json:"expirationDateTime"`
	TransactionFromDateTime string   `json:"transactionFromDateTime"`
	TransactionToDateTime   string   `json:"transactionToDateTime"`
}

type ConsentResponse struct {
	ConsentId               string   `json:"consentId"`
	Status                  string   `json:"status"`
	StatusUpdateDateTime    string   `json:"statusUpdateDateTime"`
	CreationDateTime        string   `json:"creationDateTime"`
	Permissions             []string `json:"permissions"`
	ExpirationDateTime      string   `json:"expirationDateTime"`
	TransactionFromDateTime string   `json:"transactionFromDateTime"`
	TransactionToDateTime   string   `json:"transactionToDateTime"`
}
type AccountDetails struct {
}

type Account struct {
	AccountID      string         `json:"accountId"`
	AccountDetails AccountDetails `json:"accountDetails"`
	Balance        Balance        `json:"balance"`
}

type Balance struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type Transaction struct {
	TransactionID string `json:"transactionId"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Date          string `json:"date"`
	Type          string `json:"type"`
}
