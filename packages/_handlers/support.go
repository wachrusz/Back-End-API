package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	auth "main/packages/_auth"
	email "main/packages/_email"
	jsonresponse "main/packages/_json_response"
	logger "main/packages/_logger"
)

// SupportRequest содержит информацию о запросе в техподдержку.
type SupportRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

// @Summary Send support request
// @Description Send a support request to the technical support team.
// @Tags Support
// @Accept json
// @Produce json
// @Param supportRequest body SupportRequest true "Support request object"
// @Success 200 {string} string "Support request sent successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Error sending support request"
// @Router /support/request [post]
func SendSupportRequestHandler(w http.ResponseWriter, r *http.Request) {
	var supportRequest SupportRequest
	if err := json.NewDecoder(r.Body).Decode(&supportRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	if err := sendSupportRequest(supportRequest, userID); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error sending support request: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":     "Successfully sent a suuport request",
		"status_code": http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

func sendSupportRequest(request SupportRequest, userID string) error {
	body := fmt.Sprintf("Name: %s\nEmail: %s\nSubject: %s\n\nMessage:\n%s\nUserID: %s",
		request.Name, request.Email, request.Subject, request.Message, request.UserID)

	err := email.SendEmail("support@yourdomain.com", "Support Request", body)
	if err != nil {
		logger.ErrorLogger.Printf("Error sending support request email: %v", err)
		return err
	}

	return nil
}
