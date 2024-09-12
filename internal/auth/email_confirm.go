//go:build !exclude_swagger
// +build !exclude_swagger

// Package profile provides profile information and it's functionality.
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wachrusz/Back-End-API/internal/auth/service"
	email_conf "github.com/wachrusz/Back-End-API/internal/email"
	"github.com/wachrusz/Back-End-API/internal/user"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	"github.com/wachrusz/Back-End-API/pkg/logger"
	mydb "github.com/wachrusz/Back-End-API/pkg/mydatabase"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"net/http"
	"time"
)

type ConfirmEmailRequest struct {
	Token       string `json:"token"`
	EnteredCode string `json:"code"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// ! Доделать
func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("Empty token")
		jsonresponse.SendErrorResponse(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired token: "+err.Error()), http.StatusBadRequest)
		return
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired token: "+err.Error()), http.StatusBadRequest)
		return
	}

	codeCheckResponse := email_conf.CheckConfirmationCode(registerRequest.Email, confirmRequest.Token, confirmRequest.EnteredCode)
	if codeCheckResponse.Err != "nil" {
		w.WriteHeader(codeCheckResponse.StatusCode)
		json.NewEncoder(w).Encode(codeCheckResponse)
		return
	}

	err = email_conf.ConfirmEmail(registerRequest.Email, confirmRequest.EnteredCode)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error confirming email: "+err.Error()), http.StatusInternalServerError)
		return
	}

	err = user.RegisterUser(registerRequest.Email, registerRequest.Password)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error registring user: "+err.Error()), http.StatusInternalServerError)
		return
	}
	//! SESSIONS

	userID, err_id := service.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting session: %v", err_id), http.StatusInternalServerError)
		logger.ErrorLogger.Printf("Unknown exeption in userID %s\n", userID)
	}

	deviceID, err := service.GetDeviceIDFromRequest(r)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Internal Server Error: "+err.Error()), http.StatusInternalServerError)
		return
	}
	token_details, err := generateToken(userID, deviceID, time.Minute*15)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Internal Server Error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	if service.IsUserActive(userID) {
		currentUser := service.ActiveUsers[userID]

		service.RemoveSessionFromDatabase(currentUser.DeviceID, currentUser.UserID)
		currentUser.DeviceID = deviceID
		service.ActiveUsers[userID] = currentUser
	}

	//! SAVE SESSIONS
	service.AddActiveUser(userID, registerRequest.Email, deviceID, token_details.AccessToken)

	service.SaveSessionToDatabase(registerRequest.Email, deviceID, userID, token_details.AccessToken)

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":                 "Successfuly confirmed email",
		"token_details":           token_details,
		"access_token_life_time":  time.Minute * 15,
		"refresh_token_life_time": 30 * 24 * time.Hour,
		"status_code":             http.StatusOK,
		"device_id":               deviceID,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func ConfirmEmailLoginHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("Empty token")
		jsonresponse.SendErrorResponse(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired token: "+err.Error()), http.StatusBadRequest)
		return
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired token: "+err.Error()), http.StatusBadRequest)
		return
	}

	codeCheckResponse := email_conf.CheckConfirmationCode(registerRequest.Email, confirmRequest.Token, confirmRequest.EnteredCode)
	if codeCheckResponse.Err != "nil" {
		w.WriteHeader(codeCheckResponse.StatusCode)
		json.NewEncoder(w).Encode(codeCheckResponse)
		return
	}

	err = email_conf.ConfirmEmail(registerRequest.Email, confirmRequest.EnteredCode)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error confirming email: "+err.Error()), http.StatusInternalServerError)
		return
	}

	//! SESSIONS

	userID, err_id := service.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting session: %v", err_id), http.StatusInternalServerError)
		logger.ErrorLogger.Printf("Unknown exeption in userID %s\n", userID)
	}

	deviceID, err := service.GetDeviceIDFromRequest(r)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Internal Server Error: "+err.Error()), http.StatusInternalServerError)
		return
	}
	token_details, err := generateToken(userID, deviceID, time.Minute*15)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Internal Server Error: "+err.Error()), http.StatusInternalServerError)
		return
	}

	if service.IsUserActive(userID) {
		currentUser := service.ActiveUsers[userID]

		service.RemoveSessionFromDatabase(currentUser.DeviceID, currentUser.UserID)
		currentUser.DeviceID = deviceID
		service.ActiveUsers[userID] = currentUser
	}

	//! SAVE SESSIONS
	service.AddActiveUser(userID, registerRequest.Email, deviceID, token_details.AccessToken)

	service.SaveSessionToDatabase(registerRequest.Email, deviceID, userID, token_details.AccessToken)

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":                 "Successfuly logged in",
		"token_details":           token_details,
		"access_token_life_time":  time.Minute * 15,
		"refresh_token_life_time": 30 * 24 * time.Hour,
		"status_code":             http.StatusOK,
		"device_id":               deviceID,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func ResetPasswordConfirmHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("Empty 'Content-Type' HEADER")
		jsonresponse.SendErrorResponse(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var confirmRequest ConfirmEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&confirmRequest); err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	token := confirmRequest.Token
	if token == "" {
		err := errors.New("Empty token")
		jsonresponse.SendErrorResponse(w, errors.New("Token is required: "+err.Error()), http.StatusBadRequest)
		return
	}

	claims, err := utility.ParseResetToken(token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error in parsing token: "+err.Error()), http.StatusInternalServerError)
		return
	}

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err = utility.VerifyResetJWTToken(token)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Invalid or expired token: "+err.Error()), http.StatusUnauthorized)
		return
	}

	codeCheckResponse := email_conf.CheckConfirmationCode(registerRequest.Email, confirmRequest.Token, confirmRequest.EnteredCode)
	if codeCheckResponse.Err != "nil" {
		json.NewEncoder(w).Encode(codeCheckResponse)
		return
	}

	err = email_conf.ConfirmEmail(registerRequest.Email, confirmRequest.EnteredCode)
	if err != nil {
		jsonresponse.SendErrorResponse(w, errors.New("Error confirming email: "+err.Error()), http.StatusInternalServerError)
		return
	}
	claims["confirmed"] = true

	response := map[string]interface{}{
		"message":     "Successfully confirmed email",
		"status_code": http.StatusOK,
	}
	w.WriteHeader(response["status_code"].(int))
	json.NewEncoder(w).Encode(response)
}

func ResetPassword(email, password string) error {
	hashedPassword, err := user.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = mydb.GlobalDB.Exec("UPDATE users SET hashed_password = $1 WHERE email = $2", hashedPassword, email)
	return err
}
