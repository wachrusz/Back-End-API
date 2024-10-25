package v1

import (
	"encoding/json"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/models"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
)

// AddConnectedAccountHandler handles the creation of a new connected account.
//
// @Summary Create a connected account
// @Description Create a new connected account.
// @Tags App
// @Accept json
// @Produce json
// @Param ConnectedAccount body models.ConnectedAccount true "ConnectedAccount object"
// @Success 201 {object} jsonresponse.IdResponse "Connected account created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error adding connected account"
// @Security JWT
// @Router /app/accounts [post]
func (h *MyHandler) AddConnectedAccountHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Adding a new connected account...")

	var account models.ConnectedAccount
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}
	account.UserID = userID

	connectedAccountID, err := models.AddConnectedAccount(&account)
	if err != nil {
		h.errResp(w, fmt.Errorf("error adding connected account: %v", err), http.StatusInternalServerError)
		return
	}

	response := jsonresponse.IdResponse{
		Message:    "Connected account added successfully",
		Id:         connectedAccountID,
		StatusCode: http.StatusCreated,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// DeleteConnectedAccountHandler handles the deletion of an existing connected account.
//
// @Summary Delete a connected account
// @Description Delete an existing connected account.
// @Tags App
// @Param ConnectedAccount body models.ConnectedAccount true "ConnectedAccount object"
// @Success 204 {object} jsonresponse.SuccessResponse "Connected account deleted successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error deleting connected account"
// @Security JWT
// @Router /app/accounts [delete]
func (h *MyHandler) DeleteConnectedAccountHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting connected account...")

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	err := models.DeleteConnectedAccount(userID)
	if err != nil {
		h.errResp(w, fmt.Errorf("error deleting connected account: %v", err), http.StatusInternalServerError)
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted connected account",
		StatusCode: http.StatusNoContent,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
