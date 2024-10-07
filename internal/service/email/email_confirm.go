//go:build !exclude_swagger
// +build !exclude_swagger

// Package profile provides profile information and it's functionality.
package email

import (
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	"github.com/wachrusz/Back-End-API/pkg/logger"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
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

func (s *Service) ConfirmEmail(token, code, deviceID string) (*TokenDetails, error) {
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	codeCheckResponse := s.CheckConfirmationCode(registerRequest.Email, token, code) // TODO:убрать эту хуйню
	if codeCheckResponse.Err != "nil" {
		return nil, myerrors.ErrInternal
	}

	err = s.confirmEmail(registerRequest.Email, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	err = s.users.Register(registerRequest.Email, registerRequest.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: error registring user: %v", myerrors.ErrInternal, err)
	}
	//! SESSIONS

	userID, err := s.users.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		logger.ErrorLogger.Printf("Unknown exeption in userID %s\n", userID)
		// FIXME: тут не было ретурна, но по идее должен быть
		// return fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	tokenDetails, err := s.users.GenerateToken(userID, deviceID, time.Minute*15)
	if err != nil {
		return nil, myerrors.ErrInternal
	}

	if s.users.IsUserActive(userID) {
		currentUser := s.users.ActiveUsers[userID]

		s.users.RemoveSessionFromDatabase(currentUser.DeviceID, currentUser.UserID) // TODO: handle error + concurrency safety
		currentUser.DeviceID = deviceID
		s.users.ActiveUsers[userID] = currentUser
	}

	//! SAVE SESSIONS
	s.users.AddActiveUser(userID, registerRequest.Email, deviceID, tokenDetails.AccessToken)

	s.users.SaveSessionToDatabase(registerRequest.Email, deviceID, userID, tokenDetails.AccessToken)

	return tokenDetails, nil
}

func (s *Service) ConfirmEmailLogin(token, code, deviceID string) (*user.TokenDetails, error) {
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	codeCheckResponse := s.CheckConfirmationCode(registerRequest.Email, token, code)
	if codeCheckResponse.Err != "nil" {
		return nil, myerrors.ErrInternal
	}

	err = s.confirmEmail(registerRequest.Email, code)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	//! SESSIONS

	userID, err := s.users.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		logger.ErrorLogger.Printf("Unknown exeption in userID %s\n", userID)
	}

	tokenDetails, err := s.users.GenerateToken(userID, deviceID, time.Minute*15)
	if err != nil {
		return nil, myerrors.ErrInternal
	}

	if s.users.IsUserActive(userID) {
		currentUser := s.users.ActiveUsers[userID]

		s.users.RemoveSessionFromDatabase(currentUser.DeviceID, currentUser.UserID) // TODO: handle error + concurrency safety
		currentUser.DeviceID = deviceID
		s.users.ActiveUsers[userID] = currentUser
	}

	s.users.AddActiveUser(userID, registerRequest.Email, deviceID, tokenDetails.AccessToken)

	s.users.SaveSessionToDatabase(registerRequest.Email, deviceID, userID, tokenDetails.AccessToken)

	return tokenDetails, nil
}

func (s *Service) ResetPasswordConfirm(token, code string) error {
	claims, err := utility.ParseResetToken(token)
	if err != nil {
		return myerrors.ErrInternal
	}

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err = utility.VerifyResetJWTToken(token)
	if err != nil {
		return myerrors.ErrInvalidToken
	}

	codeCheckResponse := s.CheckConfirmationCode(registerRequest.Email, token, code)
	if codeCheckResponse.Err != "nil" {
		return myerrors.ErrInternal
	}

	err = s.confirmEmail(registerRequest.Email, code)
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)

	}
	claims["confirmed"] = true
	return nil
}

func (s *Service) ResetPassword(email, password string) error {
	hashedPassword, err := s.users.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = mydb.GlobalDB.Exec("UPDATE users SET hashed_password = $1 WHERE email = $2", hashedPassword, email)
	return err
}
