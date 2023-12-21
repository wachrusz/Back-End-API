package jsonresponse

import (
	"encoding/json"
	"net/http"
)

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
