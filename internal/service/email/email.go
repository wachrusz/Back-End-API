package email

import (
	"database/sql"
	"errors"
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	"github.com/wachrusz/Back-End-API/pkg/logger"
	utility "github.com/wachrusz/Back-End-API/pkg/util"

	"net/http"
	"time"
	//"github.com/go-gomail/gomail"
)

type Service struct {
	repo *mydb.Database
}

type Emails interface {
	SendEmail(to, subject, body string) error
	SendConfirmationEmail(email, token string) error
	DecryptToken(encryptedToken string) (string, error)
	SaveConfirmationCode(email, confirmationCode, token string) error
	CheckConfirmationCode(email, token, enteredCode string) CheckResult
	DeleteConfirmationCode(email string, code string) error
	GetConfirmationCode(email string) (string, error)
	ResetPasswordConfirm(token, code string) error
	ResetPassword(email, password string) error
}

func NewService(db *mydb.Database) *Service {
	return &Service{
		repo: db,
	}
}

var (
	//corpEmail         string        = ""
	//corpEmailPassword string        = ""
	//webURL            string        = ""
	maxAttempts int = 3
	//lockDuration      time.Duration = time.Minute * 5
)

type CheckResult struct {
	RemainingAttempts int           `json:"remaining_attempts"`
	LockDuration      time.Duration `json:"lock_duration"`
	Err               string        `json:"error"`
	StatusCode        int           `json:"status_code"`
}

func (s *Service) SendEmail(to, subject, body string) error { /*
		message := mail.NewMessage()
		message.SetHeader("From", config.Username)
		message.SetHeader("To", to)
		message.SetHeader("Subject", subject)
		message.SetBody("text/plain", body)

		dialer := mail.NewDialer(config.Host, config.Port, config.Username, config.Password)

		// Send the email
		if err := dialer.DialAndSend(message); err != nil {
			return fmt.Errorf("failed to send email: %v", err)
		}
	*/
	return nil
}

func (s *Service) SendConfirmationEmail(email, token string) error {
	confirmationCode, err := utility.GenerateConfirmationCode()
	if err != nil {
		logger.ErrorLogger.Printf("Error in generating confirmation code for Email: %v", email)
		return err
	}
	/*
		m := gomail.NewMessage()
		m.SetHeader("From", corpEmail)
		m.SetHeader("To", email)
		m.SetHeader("Subject", "Код подтверждения")
		m.SetBody("text/plain", "Ваш код подтверждения: "+confirmationCode)

		d := gomail.NewDialer(webURL, 587, corpEmail, corpEmailPassword)

		if err := d.DialAndSend(m); err != nil {
			logger.ErrorLogger.Printf("Error: %v", err)
			return err
		}
	*/

	return s.SaveConfirmationCode(email, confirmationCode, token)
}

func (s *Service) DecryptToken(encryptedToken string) (string, error) {
	return enc.DecryptToken(encryptedToken)
}

func (s *Service) SaveConfirmationCode(email, confirmationCode, token string) error {
	encryptedToken, errEnc := enc.EncryptToken(token)
	if errEnc != nil {
		return errEnc
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	_, err := s.repo.Exec("INSERT INTO confirmation_codes (email, code, expiration_time, token) VALUES ($1, $2, $3, $4)", email, confirmationCode, expirationTime, encryptedToken)
	return err
}

// ! ДОДЕЛАТЬ
func (s *Service) CheckConfirmationCode(email, token, enteredCode string) CheckResult {
	var result CheckResult
	var expirationTime time.Time
	var attempts int
	result.StatusCode = http.StatusOK
	result.Err = "nil"

	err := s.checkToken(token, email)
	if err != nil {
		fmt.Println(err)
		result.Err = errors.New("Invalid token for the last code").Error()
		result.StatusCode = http.StatusUnauthorized
		return result
	}

	locked, err := s.isUserLocked(email)
	if err != nil {
		result.Err = err.Error()
		result.RemainingAttempts = 0
		result.LockDuration = time.Minute * 15
		result.StatusCode = http.StatusUnauthorized
		return result
	}
	if locked {
		result.Err = errors.New("User is locked. Try again later.").Error()
		result.RemainingAttempts = 0
		lockDuration, err := s.getLockDuration(email)
		if err != nil {
			result.Err = errors.New("Server error").Error()
		}
		result.LockDuration = lockDuration
		result.StatusCode = http.StatusUnauthorized
		return result
	}

	err = s.repo.QueryRow(`SELECT attempts FROM confirmation_codes WHERE email = $1 
	ORDER BY expiration_time DESC LIMIT 1`, email).Scan(&attempts)
	if errors.Is(err, sql.ErrNoRows) {
		result.Err = errors.New("Can't get attempts").Error()
		result.StatusCode = http.StatusUnauthorized
		return result
	} else if err != nil {
		result.Err = err.Error()
		result.StatusCode = http.StatusUnauthorized
		return result
	}

	err = s.repo.QueryRow("SELECT expiration_time FROM confirmation_codes WHERE email = $1 AND code = $2", email, enteredCode).Scan(&expirationTime)
	if errors.Is(err, sql.ErrNoRows) {
		result.Err = errors.New("Invalid confirmation code.").Error()
		result.RemainingAttempts = maxAttempts - (attempts + 1)

		if maxAttempts-(attempts+1) == 0 {
			err := s.lockUser(email)
			if err != nil {
				result.Err = err.Error()
				return result
			}

			result.Err = errors.New("User is locked. Try again later.").Error()
			lockDuration, err := s.getLockDuration(email)
			if err != nil {
				result.Err = errors.New("Server error").Error()
			}
			result.LockDuration = lockDuration
			result.StatusCode = http.StatusUnauthorized
			return result
		}

		err = s.incrementAttempts(email)
		if err != nil {
			result.Err = err.Error()
			return result
		}
		result.StatusCode = http.StatusUnauthorized
		result.RemainingAttempts = maxAttempts - (attempts + 1)
		return result
	} else if err != nil {
		result.Err = err.Error()
		return result
	}

	if time.Now().After(expirationTime) {
		result.Err = errors.New("Confirmation code has expired.").Error()
		result.StatusCode = http.StatusUnauthorized
		return result
	}

	return result
}

func (s *Service) checkToken(token, email string) error {
	var encToken string
	err := s.repo.QueryRow(`SELECT token FROM confirmation_codes WHERE email = $1
	ORDER BY expiration_time DESC LIMIT 1`, email).Scan(&encToken)
	if errors.Is(err, sql.ErrNoRows) {
		return err
	} else if err != nil {
		return err
	}

	decToken, err := s.DecryptToken(encToken)
	if err != nil {
		return err
	}

	if token != decToken {
		return errors.New("Invalid token for code")
	}
	return nil
}

func (s *Service) DeleteConfirmationCode(email string, code string) error {
	err := s.repo.QueryRow("DELETE FROM confirmation_codes WHERE email = $1 AND code = $2", email, code)
	if err != nil {
		logger.ErrorLogger.Printf("Error deleting code for Email: %v", email)
		return fmt.Errorf("error deleting confirmation")
	}
	return nil
}

func (s *Service) incrementAttempts(email string) error {
	query := `UPDATE confirmation_codes SET attempts = attempts + 1 WHERE email = $1 AND expiration_time = (
		SELECT MAX(expiration_time) 
		FROM confirmation_codes 
		WHERE email = $1
	);`
	_, err := s.repo.Exec(query, email)
	return err
}

func (s *Service) lockUser(email string) error {
	query := `UPDATE confirmation_codes SET attempts = 0, 
	locked_until = NOW() + (5 * interval '1 minute') WHERE email = $1 
		AND expiration_time = (
		SELECT MAX(expiration_time) 
		FROM confirmation_codes 
		WHERE email = $1
	);`
	_, err := s.repo.Exec(query, email)
	return err
}

func (s *Service) isUserLocked(email string) (bool, error) {
	query := "SELECT locked_until FROM confirmation_codes WHERE email = $1 ORDER BY locked_until DESC LIMIT 1"
	var lockedUntil time.Time
	err := s.repo.QueryRow(query, email).Scan(&lockedUntil)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return time.Now().Before(lockedUntil), nil
}

func (s *Service) getLockDuration(email string) (time.Duration, error) {
	var lockedUntil time.Time

	query := "SELECT locked_until FROM confirmation_codes WHERE email = $1 ORDER BY locked_until DESC LIMIT 1"
	err := s.repo.QueryRow(query, email).Scan(&lockedUntil)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return time.Until(lockedUntil), nil
}

// ! DEV/TEST functions
func (s *Service) GetConfirmationCode(email string) (string, error) {
	var code string
	err := s.repo.QueryRow("SELECT code FROM confirmation_codes WHERE email = $1 ORDER BY expiration_time DESC LIMIT 1", email).Scan(&code)
	if err != nil {
		logger.ErrorLogger.Printf("Confirmation code not found for Email: %v", email)
		return "", err
	}
	return code, nil
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

	err = s.DeleteConfirmationCode(registerRequest.Email, code)
	if err != nil {
		return fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)

	}
	claims["confirmed"] = true
	return nil
}

func (s *Service) ResetPassword(email, password string) error {
	hashedPassword, err := utility.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = s.repo.Exec("UPDATE users SET hashed_password = $1 WHERE email = $2", hashedPassword, email)
	return err
}
