package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
)

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
		_, err = s.RevokeDevice(userID, deviceID)
		if err != nil {
			return myerrors.ErrInternal
		}
		return myerrors.ErrDualSession
	}

	return err
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

// Revoke disables all user sessions on this device
func (s *Service) RevokeDevice(userID, deviceID string) (int64, error) {
	q := `
	UPDATE sessions 
	SET revoked=true 
	WHERE user_id=$1 AND device_id=$2`

	result, err := s.repo.Exec(q, userID, deviceID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Revoke disables all user sessions
func (s *Service) Revoke(userID string) (int64, error) {
	q := `
	UPDATE sessions 
	SET revoked=true 
	WHERE user_id=$1`

	result, err := s.repo.Exec(q, userID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
