package user

import (
	"database/sql"
	"errors"
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	emailConf "github.com/wachrusz/Back-End-API/internal/service/email"
	"github.com/wachrusz/Back-End-API/pkg/logger"
	utility "github.com/wachrusz/Back-End-API/pkg/util"
	"golang.org/x/crypto/bcrypt"
)

// UserAuthenticationRequest is for auth requests
type UserAuthenticationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) PrimaryRegistration(email, password string) (string, error) {
	err, used := isEmailUsed(email)
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

	err = emailConf.SendConfirmationEmail(email, token)
	if err != nil {
		return "", fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	return token, nil
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func (s *Service) ResetPassword(email string) error {
	token, err := utility.GenerateResetJWTToken(email)
	if err != nil {
		return fmt.Errorf("error generating confirmation token: %v", err)
	}

	err = emailConf.SendConfirmationEmail(email, token)
	if err != nil {
		return fmt.Errorf("error sending confirm email: %v", err)
	}

	return nil
}

type UserPasswordReset struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ResetToken string `json:"reset_token"`
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

	err = emailConf.ResetPassword(email, password)
	if err != nil {
		return fmt.Errorf("invalid email: %v", err)
	}

	userID, _ := GetUserIDFromUsersDatabase(email)
	err = s.invalidateTokensByUserID(userID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("%w: %v", myerrors.ErrInvalidToken, err)
	}

	return nil
}

// *NEW
func isEmailUsed(email string) (error, bool) {

	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	err := mydb.GlobalDB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return fmt.Errorf("Error getting email: %v", err), false
	}

	return nil, count > 0
}

type IdentificationData struct {
	Email          string
	HashedPassword string
}

func (s *Service) Register(email, password string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.getUserByEmail(email); exists {
		errMsg := "User with email " + email + " already exists"
		logger.ErrorLogger.Println(errMsg)
		return errors.New("Already exists")
	}

	if email == "" || password == "" {
		return errors.New("Blank fields are not allowed")
	}

	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.repo.Exec("INSERT INTO users (email, hashed_password) VALUES ($1, $2)", email, hashedPassword)
	if err != nil {
		logger.ErrorLogger.Println("Error inserting user:", err)
		return err
	}

	logger.InfoLogger.Printf("New user registered: %s\n", email)

	return nil
}

func (s *Service) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLogger.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *Service) getUserByEmail(email string) (IdentificationData, bool) {
	var user IdentificationData
	var id int

	row := s.repo.QueryRow("SELECT id, email, hashed_password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &user.Email, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return user, false
	} else if err != nil {
		logger.ErrorLogger.Println("Error querying user:", err)
		return user, false
	}
	return user, true
}

func (s *Service) getHashedPasswordByUsername(email string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.getUserByEmail(email)
	if !exists {
		errMsg := "User with email " + email + " not found"
		logger.ErrorLogger.Println(errMsg)
		return "", fmt.Errorf("user not found")
	}

	return user.HashedPassword, nil
}

func (s *Service) invalidateTokensByUserID(userID string) error {
	_, err := s.repo.Exec(`DELETE FROM sessions WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	return nil
}
