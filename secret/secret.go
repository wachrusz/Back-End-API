package secret

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Secret SecretValues

func init() {
	err := godotenv.Load("secret/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	Secret = SecretValues{
		CrtPath:          "secret/ok_server.crt",
		KeyPath:          "secret/ok_server.key",
		DBURL:            "postgres://postgres:" + os.Getenv("PASSWORD") + "@" + os.Getenv("HOST") + ":5432/backend_api?sslmode=disable",
		SecretKey:        []byte(os.Getenv("SECRET_KEY")),
		SecretRefreshKey: []byte(os.Getenv("SECRET_REFRESH_KEY")),
		BaseURL:          os.Getenv("HOST_VAL") + ":8080",
		CurrencyURL:      os.Getenv("CURRENCY_URL"),
		Host:             os.Getenv("HOST_VAL"),
	}
}

type SecretValues struct {
	CrtPath          string
	KeyPath          string
	DBURL            string
	SecretKey        []byte
	SecretRefreshKey []byte
	BaseURL          string
	CurrencyURL      string
	Host             string
}
