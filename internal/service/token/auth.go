package token

import (
	"database/sql"
	"fmt"

	"github.com/wachrusz/Back-End-API/internal/myerrors"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

type ConfirmEmailRequest struct {
	Token       string `json:"token"`
	EnteredCode string `json:"code"`
}

type ResetTokenDetails struct {
	ResetToken        string `json:"reset_token"`
	ExpiresAt         int64  `json:"expires_at"`
	RemainingAttempts int    `json:"-"`
	LockDuration      int    `json:"-"`
}

func (s *Service) checkLoginConds(email, password string) (bool, error) {
	if email == "" || password == "" {
		return false, myerrors.ErrEmpty
	}

	if !s.isValidCredentials(email, password) {
		return false, myerrors.ErrInvalidCreds
	}

	return true, nil
}

func (s *Service) isValidCredentials(username, password string) bool {
	hashedPassword, ok := s.getHashedPasswordByUsername(username)
	if ok != nil {
		return false
	}
	if comparePasswords(hashedPassword, password) {
		return true
	}
	return false
}

func (s *Service) getHashedPasswordByUsername(email string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	userData, exists := s.user.GetUserByEmail(email)
	if !exists {
		return "", fmt.Errorf("user not found")
	}

	return userData.HashedPassword, nil
}

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *Service) isEmailUsed(email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	err := s.repo.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error getting email: %v", err)
	}

	return count > 0, nil
}

func (s *Service) Register(email, password string) (string, error) {
	used, err := s.isEmailUsed(email)
	if err != nil {
		return "", err
	}
	if used {
		return "", myerrors.ErrDuplicated
	}

	token, err := utility.GenerateRegisterJWTToken(email, password)
	if err != nil {
		return "", fmt.Errorf("error generating confirmation token: %v", err)
	}

	err = s.email.SendConfirmationEmail(email, token)
	if err != nil {
		return "", fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	return token, nil
}

func (s *Service) Login(email, password string) (string, error) {
	if ok, err := s.checkLoginConds(email, password); !ok {
		return "", err
	}

	token, err := utility.GenerateRegisterJWTToken(email, password)
	if err != nil {
		return "", fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.email.SendConfirmationEmail(email, token)
	if err != nil {
		return "", fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	return token, nil
}

func (s *Service) ConfirmEmailRegister(token, code, deviceID string) (Details, error) {
	result := Details{}
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	codeCheckResponse, err := s.email.CheckConfirmationCode(registerRequest.Email, token, code)
	result.RemainingAttempts = codeCheckResponse.RemainingAttempts
	result.LockDuration = codeCheckResponse.LockDuration
	if err != nil {
		return result, myerrors.ErrInternal
	}

	err = s.email.DeleteConfirmationCode(registerRequest.Email, code)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	err = s.user.Register(registerRequest.Email, registerRequest.Password)
	if err != nil {
		return result, fmt.Errorf("%w: error registring user: %v", myerrors.ErrInternal, err)
	}

	userID, err := s.user.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	tokenDetails, err := s.GenerateToken(userID, deviceID)
	if err != nil {
		return result, myerrors.ErrInternal
	}

	_, err = s.RevokeDevice(userID, deviceID)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.saveSessionToDatabase(userID, deviceID, tokenDetails.RefreshToken, tokenDetails.ExpiresAt)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return tokenDetails, nil
}

func (s *Service) ConfirmEmailLogin(token, code, deviceID string) (Details, error) {
	result := Details{}
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	codeCheckResponse, err := s.email.CheckConfirmationCode(registerRequest.Email, token, code)
	result.RemainingAttempts = codeCheckResponse.RemainingAttempts
	result.LockDuration = codeCheckResponse.LockDuration
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.email.DeleteConfirmationCode(registerRequest.Email, code)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	userID, err := s.user.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	tokenDetails, err := s.GenerateToken(userID, deviceID)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	_, err = s.RevokeDevice(userID, deviceID)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.saveSessionToDatabase(userID, deviceID, tokenDetails.RefreshToken, tokenDetails.ExpiresAt)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return tokenDetails, nil
}

func (s *Service) ResetPassword(email string) (string, error) {
	exists, err := s.isEmailUsed(email)
	if err != nil {
		return "", fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	if !exists {
		return "", fmt.Errorf("%w: there is no user with this email", myerrors.ErrInvalidCreds)
	}

	token, _, err := utility.GenerateResetJWTToken(email)
	if err != nil {
		return "", fmt.Errorf("%w: error generating confirmation token: %v", myerrors.ErrInternal, err)
	}

	err = s.email.SendConfirmationEmail(email, token)
	if err != nil {
		return "", fmt.Errorf("%w: error sending confirm email: %v", myerrors.ErrEmailing, err)
	}

	return token, nil
}

func (s *Service) ResetPasswordConfirm(token, code, deviceID string) (ResetTokenDetails, error) {
	result := ResetTokenDetails{}
	email, err := utility.GetEmailFromJWT(token)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	err = utility.VerifyResetJWTToken(token, email)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	codeCheckResponse, err := s.email.CheckConfirmationCode(email, token, code)
	result.RemainingAttempts = codeCheckResponse.RemainingAttempts
	result.LockDuration = codeCheckResponse.LockDuration
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	newToken, expiresAt, err := utility.GenerateResetJWTToken(email)
	if err != nil {
		return result, fmt.Errorf("%w: error generating confirmation token: %v", myerrors.ErrInternal, err)
	}

	err = s.email.DeleteConfirmationCode(email, code)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	userID, err := s.user.GetUserIDFromUsersDatabase(email)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.saveSessionToDatabase(userID, deviceID, newToken, expiresAt)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	result.ExpiresAt = expiresAt
	result.ResetToken = newToken
	return result, nil
}

func (s *Service) ChangePasswordForRecover(email, password, resetToken string) error {
	if resetToken == "" {
		return myerrors.ErrEmpty
	}
	err := utility.VerifyResetJWTToken(resetToken, email)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token: %v", err)
	}
	claims, err := utility.ParseResetToken(resetToken)
	if claims["code_used"].(bool) {
		return fmt.Errorf("token has already been used: %v", err)
	} else {
		claims["code_used"] = true
	}

	err = s.email.ResetPassword(email, password)
	if err != nil {
		return fmt.Errorf("invalid email: %v", err)
	}

	userID, _ := s.user.GetUserIDFromUsersDatabase(email)
	_, err = s.Revoke(userID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	return nil
}

func (s *Service) Logout(device, userID string) error {
	_, err := s.RevokeDevice(device, userID)
	return err
}
