package v1

import (
	"encoding/json"
	"fmt"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"go.uber.org/zap"
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

// SendSupportRequestHandler sends a support request to the technical support team.
//
// @Summary Send support request
// @Description Send a support request to the technical support team.
// @Tags Support
// @Accept json
// @Produce json
// @Param supportRequest body SupportRequest true "Support request object"
// @Success 200 {string} string "Support request sent successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "User not authenticated"
// @Failure 500 {string} string "Error sending support request"
// @Router /support/request [post]
func (h *MyHandler) SendSupportRequestHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Sending support request...")

	// Decode the support request from the request body
	var supportRequest SupportRequest
	if err := json.NewDecoder(r.Body).Decode(&supportRequest); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	// Extract user ID from the request context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	// Convert UserID from string to int
	UserID, err := strconv.Atoi(supportRequest.UserID)
	if err != nil {
		h.l.Error("Error converting UserID", zap.Error(err))
		h.errResp(w, fmt.Errorf("invalid UserID: %v", err), http.StatusBadRequest)
		return
	}
	supportRequest.RequestID = (time.Now().UnixMicro() / 1e11) * int64(UserID)

	// Send the support request
	if err := h.sendSupportRequest(supportRequest, userID); err != nil {
		h.errResp(w, fmt.Errorf("error sending support request: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response
	response := map[string]interface{}{
		"message":           "Successfully sent a support request",
		"created_object_id": supportRequest.RequestID,
		"status_code":       http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)

	h.l.Debug("Support request sent successfully", zap.Int64("requestID", supportRequest.RequestID))
}

func (h *MyHandler) sendSupportRequest(request SupportRequest, userID string) error {
	body := fmt.Sprintf("Name: %s\nEmail: %s\nSubject: %s\n\nMessage:\n%s\nUserID: %s",
		request.Name, request.Email, request.Subject, request.Message, request.UserID)

	err := h.s.Emails.SendEmail("support@yourdomain.com", "Support Request", body)
	if err != nil {
		return err
	}

	return nil
}
