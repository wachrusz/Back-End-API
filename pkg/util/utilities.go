package utility

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var NoDeviceError = errors.New("please, fill X-Device-ID header with your device id")

type UserAuthenticationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const tokenExpirationMinutes = 15

func GenerateConfirmationCode() (string, error) {
	maxValue := big.NewInt(10000)

	randomNumber, err := rand.Int(rand.Reader, maxValue)
	if err != nil {
		return "", err
	}

	confirmationCode := fmt.Sprintf("%04d", randomNumber)

	return confirmationCode, nil
}

func GenerateRegisterJWTToken(email, password string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"password": password,
		"exp":      time.Now().Add(time.Minute * tokenExpirationMinutes).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateResetJWTToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email":     email,
		"exp":       time.Now().Add(time.Minute * tokenExpirationMinutes).Unix(),
		"code_used": false,
		"confirmed": false,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(enc.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyRegisterJWTToken(tokenString, enteredEmail, enteredPassword string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("Invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return errors.New("Email not found in token claims")
	}
	if enteredEmail != email {
		return errors.New("Email doesn't match")
	}
	password, ok := claims["password"].(string)
	if !ok {
		return errors.New("Password not found in token claims")
	}
	if enteredPassword != password {
		return errors.New("Password doesn't match")
	}

	return nil
}

func GetAuthFromJWT(tokenString string) (UserAuthenticationRequest, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return UserAuthenticationRequest{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return UserAuthenticationRequest{}, errors.New("Invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return UserAuthenticationRequest{}, errors.New("Email not found in token claims")
	}
	password, ok := claims["password"].(string)
	if !ok {
		return UserAuthenticationRequest{}, errors.New("Password not found in token claims")
	}

	return UserAuthenticationRequest{
			Email:    email,
			Password: password,
		},
		nil
}

func GetEmailFromJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("Email not found in token claims")
	}

	return email, nil
}

func VerifyResetJWTToken(tokenString, enteredEmail string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("Invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return errors.New("Email not found in token claims")
	}
	if enteredEmail != email {
		return errors.New("Email doesn't match")
	}

	return nil
}

func ParseResetToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return enc.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid token")
}

func GetDeviceIDFromJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(enc.SecretKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid token")
	}

	deviceID, ok := claims["device_id"].(string)
	if !ok {
		return "", errors.New("Invalid token")
	}
	return deviceID, nil
}

func GetDeviceIDFromRequest(r *http.Request) (string, error) {
	deviceID := strings.TrimSpace(r.Header.Get("X-Device-ID"))
	if deviceID == "" {
		return "", NoDeviceError
	}
	return deviceID, nil
}

// ExtractTokenFromHeader извлекает токен из заголовка Authorization
func ExtractTokenFromHeader(r *http.Request) (string, error) {
	// Получаем значение заголовка Authorization
	tokenString := r.Header.Get("Authorization")

	// Проверяем, что заголовок не пустой
	if tokenString == "" {
		return "", errors.New("missing Authorization header")
	}

	// Проверяем, что заголовок начинается с "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(tokenString, bearerPrefix) {
		return "", errors.New("invalid Authorization header format, expected 'Bearer <token>'")
	}

	// Убираем префикс "Bearer " и возвращаем токен
	token := strings.TrimSpace(strings.TrimPrefix(tokenString, bearerPrefix))
	if token == "" {
		return "", errors.New("empty token in Authorization header")
	}

	return token, nil
}

func GetDeviceIDFromContext(ctx context.Context) (string, bool) {
	deviceID, ok := ctx.Value("device_id").(string)
	return deviceID, ok
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("userID").(string)
	return userID, ok
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GetParamFromRequest extracts a parameter from the URL by its key.
func GetParamFromRequest(r *http.Request, key string) string {
	// chi.URLParam позволяет извлечь параметр из URL
	return chi.URLParam(r, key)
}
