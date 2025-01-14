package token

import (
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/internal/service/user"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
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
	Register(email, password string) (string, error)
	Login(email, password string) (string, error)
	ConfirmEmailRegister(token, code, deviceID string) (Details, error)
	ConfirmEmailLogin(token, code, deviceID string) (Details, error)
	Logout(device, userID string) error
	ResetPassword(email string) (string, error)
	ResetPasswordConfirm(token, code, deviceId string) (ResetTokenDetails, error)
	ChangePasswordForRecover(email, password, resetToken string) error
	RefreshToken(refreshTokenString, deviceID string) (Details, error)
	GenerateToken(userID string, deviceID string) (Details, error)
	RevokeDevice(userID, deviceID string) (int64, error)
	Revoke(userID string) (int64, error)
}

type Details struct {
	AccessToken       string `json:"access_token"`
	RefreshToken      string `json:"refresh_token"`
	ExpiresAt         int64  `json:"expires_at"`
	RemainingAttempts int    `json:"-"`
	LockDuration      int    `json:"-"`
}

func (s *Service) invalidateTokensByUserID(userID string) error {
	// TODO: cron
	_, err := s.repo.Exec(`DELETE FROM sessions WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	return nil
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
