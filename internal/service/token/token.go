package token

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

const (
	refreshTokenLifeTime = time.Hour * 24 * 30 // 30 days
	resetTokenLifeTime   = time.Minute * 15    // 15 minutes
)

type Service struct {
	email               email.Emails
	user                user.Users
	repo                *mydb.Database
	mutex               sync.Mutex
	accessTokenLifetime time.Duration
}

func NewService(repo *mydb.Database, e email.Emails, u user.Users, d int) *Service {
	return &Service{
		email:               e,
		user:                u,
		repo:                repo,
		mutex:               sync.Mutex{},
		accessTokenLifetime: time.Minute * time.Duration(d),
	}
}

type Tokens interface {
	PrimaryRegistration(email, password string) (string, error)
	ResetPassword(email string) (string, error)
	ResetPasswordConfirm(token, code, deviceId string) (ResetTokenDetails, error)
	ChangePasswordForRecover(email, password, resetToken string) error
	Login(email, password string) (string, error)
	RefreshToken(refreshTokenString, deviceID string) (Details, error)
	GenerateToken(userID string, deviceID string) (Details, error)
	ConfirmEmailRegister(token, code, deviceID string) (Details, error)
	ConfirmEmailLogin(token, code, deviceID string) (Details, error)
	RemoveSessionFromDatabase(deviceID, userID string) error
	Logout(device, userID string) error
}

func (s *Service) PrimaryRegistration(email, password string) (string, error) {
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

func (s *Service) isEmailUsed(email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	err := s.repo.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error getting email: %v", err)
	}

	return count > 0, nil
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
	err = s.invalidateTokensByUserID(userID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	return nil
}

func (s *Service) invalidateTokensByUserID(userID string) error {
	_, err := s.repo.Exec(`DELETE FROM sessions WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	return nil
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

func (s *Service) RefreshToken(refreshTokenString, deviceID string) (Details, error) {
	details := Details{}
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return details, err
	}

	if _, ok := refreshToken.Claims.(jwt.Claims); !ok && !refreshToken.Valid {
		return details, fmt.Errorf("Invalid refresh token")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return details, fmt.Errorf("Failed to parse refresh token claims")
	}

	userIDClaims, ok := claims["sub"].(string)
	if !ok {
		return details, fmt.Errorf("Failed to get user ID from refresh token")
	}

	deviceIDClaims, ok := claims["device_id"].(string)
	if !ok {
		return details, fmt.Errorf("Failed to get device ID from refresh token")
	}

	if deviceID != deviceIDClaims {
		return details, myerrors.ErrInvalidToken
	}

	expiresFloat, ok := claims["exp"].(float64)
	if !ok {
		return details, fmt.Errorf("Failed to get expiration time from refresh token")
	}
	expires := int64(expiresFloat)

	if time.Unix(expires, 0).Before(time.Now()) {
		return details, myerrors.ErrExpiredToken
	}

	details, err = s.GenerateToken(userIDClaims, deviceID)
	if err != nil {
		return details, err
	}

	err = s.updateTokenInDB(userIDClaims, deviceID, details.RefreshToken, details.ExpiresAt)
	if err != nil {
		return details, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return details, nil
}

// updateTokenInDB выполняет ротацию refresh-токенов и обновляет последнюю активность пользователя
func (s *Service) updateTokenInDB(userID, deviceID, token string, expirationTime int64) error {
	encryptedToken, err := enc.EncryptToken(token)
	if err != nil {
		return err
	}

	expirationTimeFormatted := time.Unix(expirationTime, 0).UTC()

	result, err := s.repo.Exec(`
        UPDATE sessions 
        SET last_activity=NOW(), token=$1, expires_at=$2
        WHERE (user_id=$3 AND device_id=$4) AND (expires_at>NOW() AND revoked=false)`,
		encryptedToken, expirationTimeFormatted, userID, deviceID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return myerrors.ErrExpiredToken
	}

	if rowsAffected != 1 {
		// TODO: revoke
		return myerrors.ErrDualSession
	}

	return err
}

func (s *Service) GenerateToken(userID string, deviceID string) (Details, error) {
	details := Details{}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"exp":       time.Now().Add(s.accessTokenLifetime).Unix(),
		"device_id": deviceID,
	})

	accessTokenString, err := accessToken.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return details, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"exp":       time.Now().Add(refreshTokenLifeTime).Unix(),
		"device_id": deviceID,
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return details, err
	}

	refreshTokenExpiresAt := time.Now().Add(refreshTokenLifeTime).Unix()

	return Details{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    refreshTokenExpiresAt,
	}, nil
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
	//! SESSIONS

	userID, err := s.user.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	tokenDetails, err := s.GenerateToken(userID, deviceID)
	if err != nil {
		return result, myerrors.ErrInternal
	}

	//! SAVE SESSIONS
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

	err = s.saveSessionToDatabase(userID, deviceID, tokenDetails.RefreshToken, tokenDetails.ExpiresAt)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return tokenDetails, nil
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

type ConfirmEmailRequest struct {
	Token       string `json:"token"`
	EnteredCode string `json:"code"`
}

type Details struct {
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	ExpiresAt         int64  `json:"expires_at"`
	RemainingAttempts int    `json:"-"`
	LockDuration      int    `json:"-"`
}

type ResetTokenDetails struct {
	ResetToken        string `json:"reset_token"`
	ExpiresAt         int64  `json:"expires_at"`
	RemainingAttempts int    `json:"-"`
	LockDuration      int    `json:"-"`
}

func (s *Service) saveSessionToDatabase(userID, deviceID, token string, expirationTime int64) error {
	// Encrypt the token
	encryptedToken, err := enc.EncryptToken(token)
	if err != nil {
		return err
	}

	expirationTimeFormatted := time.Unix(expirationTime, 0).UTC()

	_, err = s.repo.Exec(`
        INSERT INTO sessions (device_id, created_at, last_activity, user_id, token, expires_at)
        VALUES ($1, NOW(), NOW(), $2, $3, $4)`,
		deviceID, userID, encryptedToken, expirationTimeFormatted)
	return err
}

func (s *Service) RemoveSessionFromDatabase(deviceID, userID string) error {
	_, err := s.repo.Exec(`
        UPDATE sessions SET revoked = true WHERE device_id = $1 AND user_id = $2;`, deviceID, userID)
	return err
}

func (s *Service) Logout(device, userID string) error {
	err := s.RemoveSessionFromDatabase(device, userID)
	if err != nil {
		return fmt.Errorf("error removing session from db: %v", err)
	}

	return nil
}
