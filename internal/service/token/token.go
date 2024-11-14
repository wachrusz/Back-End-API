package token

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

type Service struct {
	email email.Emails
	user  user.Users
	repo  *mydb.Database
	mutex sync.Mutex
}

func NewService(repo *mydb.Database, e email.Emails, u user.Users) *Service {
	return &Service{
		email: e,
		user:  u,
		repo:  repo,
		mutex: sync.Mutex{},
	}
}

type Tokens interface {
	PrimaryRegistration(email, password string) (string, error)
	ResetPassword(email string) error
	ChangePasswordForRecover(email, password, resetToken string) error
	Login(email, password string) (string, error)
	RefreshToken(rt, userID string) (string, string, error)
	GenerateToken(userID string, deviceID string, duration time.Duration) (*Details, error)
	ConfirmEmailRegister(token, code, deviceID string) (*Details, int, int, error)
	ConfirmEmailLogin(token, code, deviceID string) (*Details, int, int, error)
}

func (s *Service) PrimaryRegistration(email, password string) (string, error) {
	err, used := s.isEmailUsed(email)
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

func (s *Service) isEmailUsed(email string) (error, bool) {

	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	err := s.repo.QueryRow(query, email).Scan(&count)
	if err != nil {
		return fmt.Errorf("error getting email: %v", err), false
	}

	return nil, count > 0
}

func (s *Service) ResetPassword(email string) error {
	token, err := utility.GenerateResetJWTToken(email)
	if err != nil {
		return fmt.Errorf("error generating confirmation token: %v", err)
	}

	err = s.email.SendConfirmationEmail(email, token)
	if err != nil {
		return fmt.Errorf("error sending confirm email: %v", err)
	}

	return nil
}

func (s *Service) ChangePasswordForRecover(email, password, resetToken string) error {
	if resetToken == "" {
		return myerrors.ErrEmpty
	}
	_, err := utility.VerifyResetJWTToken(resetToken)
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

func (s *Service) RefreshToken(rt, userID string) (string, string, error) {
	tokenDetails, err := s.refreshToken(rt)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.updateTokenInDB(userID, tokenDetails.AccessToken)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return tokenDetails.AccessToken, tokenDetails.RefreshToken, nil
}

// *NEW
// ! Это выглядит подозрительно плохо, есть вероятность что нахуй буду послан при попытке запустить(так и было)
func (s *Service) updateTokenInDB(userID, newAccessToken string) error {
	encryptedToken, err := enc.EncryptToken(newAccessToken)
	if err != nil {
		return err
	}

	s.user.SetAccessToken(userID, newAccessToken)

	query := `
		UPDATE sessions
		SET token = $1,
		expires_at = NOW() + INTERVAL '15 minutes'
		WHERE user_id = $2;
	`
	_, err = s.repo.Exec(query, encryptedToken, userID)
	return err
}

func (s *Service) refreshToken(refreshTokenString string) (*Details, error) {
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := refreshToken.Claims.(jwt.Claims); !ok && !refreshToken.Valid {
		return nil, fmt.Errorf("Invalid refresh token")
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Failed to parse refresh token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("Failed to get user ID from refresh token")
	}
	deviceID, ok := claims["device_id"].(string)
	if !ok {
		return nil, fmt.Errorf("Failed to get device ID from refresh token")
	}

	return s.GenerateToken(userID, deviceID, time.Minute*15)
}

func (s *Service) GenerateToken(userID string, device_id string, duration time.Duration) (*Details, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"exp":       time.Now().Add(duration).Unix(),
		"device_id": device_id,
	})

	accessTokenString, err := accessToken.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       userID,
		"device_id": device_id,
	})
	refreshTokenString, err := refreshToken.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return nil, err
	}

	refreshTokenExpiresAt := time.Now().Add(30 * 24 * time.Hour).Unix()

	return &Details{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    refreshTokenExpiresAt,
	}, nil
}

func (s *Service) ConfirmEmailRegister(token, code, deviceID string) (*Details, int, int, error) {
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	codeCheckResponse, err := s.email.CheckConfirmationCode(registerRequest.Email, token, code)
	if err != nil {
		return nil, codeCheckResponse.RemainingAttempts, codeCheckResponse.LockDuration, myerrors.ErrInternal
	}

	err = s.email.DeleteConfirmationCode(registerRequest.Email, code)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	err = s.user.Register(registerRequest.Email, registerRequest.Password)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: error registring user: %v", myerrors.ErrInternal, err)
	}
	//! SESSIONS

	userID, err := s.user.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	tokenDetails, err := s.GenerateToken(userID, deviceID, time.Minute*15)
	if err != nil {
		return nil, 0, 0, myerrors.ErrInternal
	}

	//! SAVE SESSIONS
	err = s.user.SaveSessionToDatabase(registerRequest.Email, deviceID, userID, tokenDetails.AccessToken)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return tokenDetails, 0, 0, nil
}

func (s *Service) ConfirmEmailLogin(token, code, deviceID string) (*Details, int, int, error) {
	registerRequest, err := utility.GetAuthFromJWT(token)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}
	err = utility.VerifyRegisterJWTToken(token, registerRequest.Email, registerRequest.Password)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	codeCheckResponse, err := s.email.CheckConfirmationCode(registerRequest.Email, token, code)
	if err != nil {
		return nil, codeCheckResponse.RemainingAttempts, codeCheckResponse.LockDuration, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.email.DeleteConfirmationCode(registerRequest.Email, code)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	//! SESSIONS
	userID, err := s.user.GetUserIDFromUsersDatabase(registerRequest.Email)
	if err != nil {

	}

	tokenDetails, err := s.GenerateToken(userID, deviceID, time.Minute*15)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.user.SaveSessionToDatabase(registerRequest.Email, deviceID, userID, tokenDetails.AccessToken)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return tokenDetails, 0, 0, nil
}

type ConfirmEmailRequest struct {
	Token       string `json:"token"`
	EnteredCode string `json:"code"`
}

type Details struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
