package jsonresponse

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

func SendErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	errorResponse := ErrorResponse{
		Error:      err.Error(),
		StatusCode: statusCode,
	}

	jsonData, jsonErr := json.Marshal(errorResponse)
	if jsonErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

type IdResponse struct {
	Message    string `json:"message"`
	Id         int64  `json:"id"`
	StatusCode int    `json:"status_code"`
}

type TokenResponse struct {
	Message    string `json:"message"`
	Token      string `json:"token"`
	StatusCode int    `json:"status_code"`
}

type DoubleTokenResponse struct {
	Message              string `json:"message"`
	AccessToken          string `json:"access_token"`
	RefreshToken         string `json:"refresh_token"`
	AccessTokenLifeTime  int64  `json:"access_token_life_time"`
	RefreshTokenLifeTime int64  `json:"refresh_token_life_time"`
	StatusCode           int    `json:"status_code"`
}

type AmountResponse struct {
	Message    string `json:"message"`
	Amount     int    `json:"amount"`
	StatusCode int    `json:"status_code"`
}

type CodeResponse struct {
	Message    string `json:"message"`
	Code       string `json:"code"`
	StatusCode int    `json:"status_code"`
}

type CodeError struct {
	Error        string `json:"error"`
	Attempts     int    `json:"remaining_attempts"`
	LockDuration int    `json:"lock_duration"`
	StatusCode   int    `json:"status_code"`
}
