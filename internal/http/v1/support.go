package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/logger"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

// SupportRequest содержит информацию о запросе в техподдержку.
type SupportRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	UserID    string `json:"user_id"`
	RequestID int64  `json:"request_id"`
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
func (h *MyHandler) SendSupportRequestHandler(w http.ResponseWriter, r *http.Request) {
	var supportRequest SupportRequest
	if err := json.NewDecoder(r.Body).Decode(&supportRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		jsonresponse.SendErrorResponse(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
		return
	}

	UserID, err := strconv.Atoi(supportRequest.UserID)
	if err != nil {
		log.Println(err)
		return
	}
	supportRequest.RequestID = (time.Now().UnixMicro() / 1e11) * int64(UserID)

	if err := h.sendSupportRequest(supportRequest, userID); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error sending support request: "+err.Error()), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":           "Successfully sent a support request",
		"created_object_id": supportRequest.RequestID,
		"status_code":       http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func (h *MyHandler) sendSupportRequest(request SupportRequest, userID string) error {
	body := fmt.Sprintf("Name: %s\nEmail: %s\nSubject: %s\n\nMessage:\n%s\nUserID: %s",
		request.Name, request.Email, request.Subject, request.Message, request.UserID)

	err := h.s.Emails.SendEmail("support@yourdomain.com", "Support Request", body)
	if err != nil {
		logger.ErrorLogger.Printf("Error sending support request email: %v", err)
		return err
	}

	return nil
}
