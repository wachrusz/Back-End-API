package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	jsonresponse "github.com/wachrusz/Back-End-API/pkg/json_response"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"github.com/wachrusz/Back-End-API/pkg/validator"
	"go.uber.org/zap"
)

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenHandler handles the refresh of authentication tokens.
//
// @Summary Refresh authentication tokens
// @Description This endpoint allows users to refresh their access and refresh tokens using an existing refresh RefreshToken.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   token  body   RefreshToken   true  "Refresh Token"
// @Success 200 {object} jsonresponse.DoubleTokenResponse "Successfully refreshed tokens"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid Content-Type or request payload"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error"
// @Security JWT
// @Router /auth/refresh [post]
func (h *MyHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		h.errResp(w, errors.New("Invalid Content-Type, expected application/json: "), http.StatusBadRequest)
		return
	}
	//! Заставляет задуматься
	/*
		userID, ok := utility.GetUserIDFromContext(r.Context())
		if !ok {
			h.errResp(w, errors.New("User not authenticated: "), http.StatusUnauthorized)
			return
		}
	*/
	var tmpToken RefreshToken

	err := json.NewDecoder(r.Body).Decode(&tmpToken)
	if err != nil {
		h.errResp(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	refreshToken, err := jwt.Parse(tmpToken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(enc.SecretKey), nil
	})

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		h.errResp(w, errors.New("no claims in token: "+err.Error()), http.StatusBadRequest)
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		h.errResp(w, errors.New("Failed to get user ID from refresh token: "+err.Error()), http.StatusBadRequest)
		return
	}

	access, refresh, err := h.s.Tokens.RefreshToken(tmpToken.RefreshToken, userID)
	if err != nil {
		h.errResp(w, err, http.StatusInternalServerError)
		return
	}

	response := jsonresponse.DoubleTokenResponse{
		Message:              "Successfully refreshed tokens",
		AccessToken:          access,
		RefreshToken:         refresh,
		AccessTokenLifeTime:  60 * 15,
		RefreshTokenLifeTime: 30 * 24 * 60 * 60,
		StatusCode:           http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// LogoutUserHandler handles user logout and terminates the session.
//
// @Summary LogoutUserHandler the user
// @Description Logs out the user, terminating the session.
// @Tags Auth
// @Produce json
// @Success 200 {object} jsonresponse.SuccessResponse "LogoutUserHandler successful"
// @Failure 401 {object} jsonresponse.ErrorResponse "User not authenticated"
// @Failure 500 {object} jsonresponse.ErrorResponse "Internal Server Error"
// @Security JWT
// @Router /auth/logout [post]
func (h *MyHandler) LogoutUserHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("User logout initiated...")

	// Retrieve the current device ID from the context
	currentDeviceID, ok := utility.GetDeviceIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated: device ID not found"), http.StatusUnauthorized)
		return
	}

	// Retrieve the user ID from the context
	userID, ok := utility.GetUserIDFromContext(r.Context())
	if !ok {
		h.errResp(w, fmt.Errorf("user not authenticated: user ID not found"), http.StatusUnauthorized)
		return
	}

	// Perform logout operation
	if err := h.s.Users.Logout(currentDeviceID, userID); err != nil {
		h.errResp(w, fmt.Errorf("error during logout: %v", err), http.StatusInternalServerError)
		return
	}

	// Send successful logout response
	response := jsonresponse.SuccessResponse{
		Message:    "Logout successful",
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

	h.l.Debug("User logged out successfully", zap.String("userID", userID))
}

// DeleteTokensHandler handles the deletion of authentication tokens.
//
// @Summary Delete authentication tokens
// @Description This endpoint allows users to delete authentication tokens either by email or device ID. Only one of the parameters should be provided.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   email    query  string  false  "User email"
// @Param   deviceID query  string  false  "Device ID"
// @Success 204 {object} jsonresponse.SuccessResponse "Successfully deleted tokens"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request: blank fields or both email and deviceID provided"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error"
// @Router /auth/tokens [delete]
func (h *MyHandler) DeleteTokensHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Deleting tokens...")
	email := r.URL.Query().Get("email")
	deviceID := r.URL.Query().Get("deviceID")
	if (email == "" && deviceID == "") || (email != "" && deviceID != "") {
		h.errResp(w, errors.New("blank fields and two methods are not allowed"), http.StatusBadRequest)
		return
	}

	err := h.s.Users.DeleteTokens(email, deviceID)
	if err != nil {
		h.errResp(w, err, http.StatusInternalServerError)
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully deleted tokens",
		StatusCode: http.StatusNoContent,
	}

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// GetTokenPairsAmountHandler retrieves the number of token pairs for a user.
//
// @Summary Get token pairs amount
// @Description This endpoint returns the number of token pairs (active sessions) associated with the provided email.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   email  query  string  true  "User email"
// @Success 200 {object} jsonresponse.AmountResponse "Successfully got the amount of token pairs"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request: blank fields are not allowed"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error"
// @Router /auth/tokens/amount [get]
func (h *MyHandler) GetTokenPairsAmountHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Getting pairs amount")
	email := r.URL.Query().Get("email")
	if email == "" {
		h.errResp(w, fmt.Errorf("blank fields are not allowed"), http.StatusBadRequest)
		return
	}
	amount, err := h.s.Users.GetTokenPairsAmount(email)
	if err != nil {
		h.errResp(w, fmt.Errorf("error while counting sessions: %v", err.Error()), http.StatusInternalServerError)
		return
	}
	response := jsonresponse.AmountResponse{
		Message:    "Successfully got amount",
		Amount:     amount,
		StatusCode: http.StatusOK,
	}
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

// ChangePasswordForRecoverHandler allows users to reset their password using a reset token.
//
// @Summary Change password for recovery
// @Description This endpoint enables users to change their password by providing their email, new password, and a valid reset token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   body  body  UserPasswordReset  true  "Password reset request with email, new password, and reset token"
// @Success 200 {object} jsonresponse.SuccessResponse "Password reset successfully"
// @Failure 400 {object} jsonresponse.ErrorResponse "Invalid request: empty fields, invalid email, or password too short"
// @Failure 401 {object} jsonresponse.ErrorResponse "Invalid or expired reset token"
// @Failure 500 {object} jsonresponse.ErrorResponse "Server error"
// @Router /auth/password [put]
func (h *MyHandler) ChangePasswordForRecoverHandler(w http.ResponseWriter, r *http.Request) {
	h.l.Debug("Changing password...")
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		err := errors.New("empty 'Content-Type' HEADER")
		h.errResp(w, errors.New("Invalid Content-Type, expected application/json: "+err.Error()), http.StatusBadRequest)
		return
	}

	var resetRequest UserPasswordReset

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		h.errResp(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
		return
	}

	email := resetRequest.Email
	password := resetRequest.Password
	if !validator.IsValidPassword(password) {
		h.errResp(w, errors.New("password must be at least 7 digits long: "), http.StatusBadRequest)
		return
	}

	if !validator.IsValidEmail(email) {
		h.errResp(w, errors.New("Invalid email: "), http.StatusBadRequest)
		return
	}

	resetToken := resetRequest.ResetToken
	err = h.s.Tokens.ChangePasswordForRecover(email, password, resetToken)
	if err != nil {
		var statusCode = 500
		switch {
		case errors.Is(err, myerrors.ErrEmpty):
			statusCode = http.StatusBadRequest
			break
		case errors.Is(err, myerrors.ErrEmailing):
			statusCode = http.StatusInternalServerError
			break
		case errors.Is(err, myerrors.ErrInvalidToken):
			statusCode = http.StatusUnauthorized
			break
		}
		h.errResp(w, fmt.Errorf("error changing the password: %v", err), statusCode)
		return
	}

	response := jsonresponse.SuccessResponse{
		Message:    "Successfully reset password",
		StatusCode: http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type UserPasswordReset struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ResetToken string `json:"reset_token"`
}
