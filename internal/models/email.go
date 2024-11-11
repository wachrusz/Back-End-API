package models

type Email struct {
	To      string `json:"to"`      // Адрес получателя
	Subject string `json:"subject"` // Тема письма
	Body    string `json:"body"`    // Тело письма
}
