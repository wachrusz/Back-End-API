package currency

import (
	"encoding/json"
	"fmt"
	"io"
	mydb "main/packages/_mydatabase"
	"main/secret"
	"net/http"
	"time"
)

var (
	CurrentCurrencyData CurrencyData
)

type Valute struct {
	ID       string  `json:"ID"`
	NumCode  string  `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Nominal  int     `json:"Nominal"`
	Name     string  `json:"Name"`
	Value    float64 `json:"Value"`
	Previous float64 `json:"Previous"`
}

type CurrencyData struct {
	Date         string `json:"Date"`
	PreviousDate string `json:"PreviousDate"`
	PreviousURL  string `json:"PreviousURL"`
	Timestamp    string `json:"Timestamp"`
	Valute       map[string]Valute
}

func InitCurrentCurrencyData() error {
	err := parseJSONAndUpdateDB(secret.Secret.CurrencyURL)
	if err != nil {
		fmt.Println("Error in updating database:", err)
	}
	CurrentCurrencyData.Valute = make(map[string]Valute)
	rows, err := mydb.GlobalDB.Query("SELECT id, num_code, currency_code, nominal, name, value, previous FROM currency")
	if err != nil {
		return err
	}

	for rows.Next() {
		var item Valute
		err := rows.Scan(&item.ID, &item.NumCode, &item.CharCode, &item.Nominal, &item.Name, &item.Value, &item.Previous)
		if err != nil {
			return err
		}

		CurrentCurrencyData.Valute[item.CharCode] = item
	}

	return nil
}

func parseJSONAndUpdateDB(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error retrieving data: %w", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Ошибка при чтении тела ответа: %w", err)
	}

	var data CurrencyData
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return fmt.Errorf("Ошибка при разборе JSON: %w", err)
	}

	CurrentCurrencyData = data

	err = updateCurrencyRatesAndDataInDB(data.Valute)
	if err != nil {
		return fmt.Errorf("Ошибка при обновлении курсов валют в базе данных: %w", err)
	}

	return nil
}

func updateCurrencyRatesAndDataInDB(rates map[string]Valute) error {
	query1 := `
    INSERT INTO exchange_rates (currency_code, rate_to_ruble)
    VALUES ($1, $2)
    ON CONFLICT (currency_code) DO UPDATE SET rate_to_ruble = EXCLUDED.rate_to_ruble;
    `

	query2 := `
    INSERT INTO currency (cbr_id, num_code, currency_code, nominal, name, value, previous)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT (currency_code) DO UPDATE SET num_code = EXCLUDED.num_code, currency_code = EXCLUDED.currency_code, nominal = EXCLUDED.nominal, name = EXCLUDED.name, value = EXCLUDED.value, previous = EXCLUDED.previous;
    `

	for _, item := range rates {
		_, err := mydb.GlobalDB.Exec(query1, item.CharCode, item.Value/float64(item.Nominal))
		if err != nil {
			return err
		}

		_, err = mydb.GlobalDB.Exec(query2, item.ID, item.NumCode, item.CharCode, item.Nominal, item.Name, item.Value, item.Nominal)
		if err != nil {
			return err
		}
	}

	return nil
}

func ScheduleCurrencyUpdates() {
	timeStr := "11:40"
	updateHour, err := time.Parse("15:04", timeStr)
	if err != nil {
		fmt.Println("Parsing error:", err)
		return
	}

	for {
		now := time.Now()
		nextUpdate := time.Date(now.Year(), now.Month(), now.Day(), updateHour.Hour(), updateHour.Minute(), 0, 0, time.Local)

		if nextUpdate.Before(now) {
			nextUpdate = nextUpdate.Add(24 * time.Hour)
		}

		time.Sleep(nextUpdate.Sub(now))

		err := parseJSONAndUpdateDB(secret.Secret.CurrencyURL)
		if err != nil {
			fmt.Println("Error in updating database:", err)
		}
	}
}
