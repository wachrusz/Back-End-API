package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/repository/models"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
)

type ConnectedAccountRequest struct {
	Account models.ConnectedAccount `json:"account"`
}

// AddConnectedAccountHandler handles the creation of a new connected account.
//
// @Summary Create a connected account
// @Description Create a new connected account.
// @Tags App
// @Accept json
// @Produce json
// @Param ConnectedAccount body ConnectedAccountRequest true "ConnectedAccount object"
// @Success 201 {object} jsonresponse.IdResponse "Connected account created successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error adding connected account"
// @Security JWT
// @Router /app/accounts [post]
func (h *MyHandler) AddConnectedAccountHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Adding a new connected account...")

	var accountR ConnectedAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&accountR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	account := accountR.Account

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}
	account.UserID = userID

	connectedAccountID, err := h.m.Accounts.Create(&account)
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
// @Param ConnectedAccount body jsonresponse.IdRequest true "ConnectedAccount id"
// @Success 204 {object} jsonresponse.SuccessResponse "Connected account deleted successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error deleting connected account"
// @Security JWT
// @Router /app/accounts [delete]
func (h *MyHandler) DeleteConnectedAccountHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting connected account...")

	var id jsonresponse.IdRequest
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}

	if err := h.m.Accounts.Delete(id.ID, userID); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("connected account not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error deleting connected account: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted connected account",
		StatusCode: http.StatusNoContent,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// UpdateConnectedAccountHandler handles the update of an existing connected account.
//
// @Summary Update a connected account
// @Description Update an existing connected account. There is no need to fill user_id field.
// @Tags App
// @Accept json
// @Produce json
// @Param ConnectedAccount body ConnectedAccountRequest true "ConnectedAccount object"
// @Success 200 {object} jsonresponse.SuccessResponse "Connected account updated successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 404 {object} jsonresponse.ErrorResponse "Connected account not found"
// @Failure 500 {object} jsonresponse.ErrorResponse "Error updating connected account"
// @Security JWT
// @Router /app/accounts [put]
func (h *MyHandler) UpdateConnectedAccountHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Updating connected account...")

	// Decode the request body
	var accountR ConnectedAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&accountR); err != nil {
		h.errResp(w, fmt.Errorf("invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	editedAccount := accountR.Account

	// Check if user is authenticated
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated"), http.StatusUnauthorized)
		return
	}
	editedAccount.UserID = userID

	// Attempt to update the account
	if err := h.m.Accounts.Update(&editedAccount); err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			h.errResp(w, fmt.Errorf("connected account not found: %v", err), http.StatusNotFound)
		} else {
			h.errResp(w, fmt.Errorf("error updating connected account: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	response := jsonresponse.SuccessResponse{
		Message:    "Connected account updated successfully",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
