//go:build !exclude_swagger
// +build !exclude_swagger

// Package auth provides authentication and authorization functionality.
package user

import (
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	"github.com/wachrusz/Back-End-API/internal/service/email"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func (s *Service) Login(email, password string) (string, error) {
	if ok, err := s.checkLoginConds(email, password); !ok {
		return "", err
	}

	token, err := utility.GenerateRegisterJWTToken(email, password)
	if err != nil {
		return "", myerrors.ErrInternal
	}

	err = email.SendConfirmationEmail(email, token)
	if err != nil {
		return "", myerrors.ErrEmailing
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

// *NEW
// ! СРОЧНО ДОДЕЛАТЬ ЭТО
func (s *Service) RefreshToken(rt, userID string) (string, string, error) {
	tokenDetails, err := s.refreshToken(rt)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = updateTokenInDB(userID, tokenDetails.AccessToken)
	if err != nil {
		return "", "", fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	return tokenDetails.AccessToken, tokenDetails.RefreshToken, nil
}

// *NEW
// ! Это выглядит подозрительно плохо, есть вероятность что нахуй буду послан при попытке запустить(так и было)
func updateTokenInDB(userID, newAccessToken string) error {
	encryptedToken, err := enc.EncryptToken(newAccessToken)
	if err != nil {
		return err
	}

	SetAccessToken(userID, newAccessToken)

	query := `
		UPDATE sessions
		SET token = $1,
		expires_at = NOW() + INTERVAL '15 minutes'
		WHERE user_id = $2;
	`
	_, err = mydb.GlobalDB.Exec(query, encryptedToken, userID)
	return err
}

func (s *Service) GenerateToken(userID string, device_id string, duration time.Duration) (*email.TokenDetails, error) {
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

	return &email.TokenDetails{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    refreshTokenExpiresAt,
	}, nil
}

func (s *Service) refreshToken(refreshTokenString string) (*email.TokenDetails, error) {
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

func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
