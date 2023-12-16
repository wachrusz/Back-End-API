package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "backEndAPI/_auth"
	email "backEndAPI/_email"
	logger "backEndAPI/_logger"
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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	if err := sendSupportRequest(supportRequest, userID); err != nil {
		http.Error(w, "Error sending support request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Support request sent successfully"))
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
