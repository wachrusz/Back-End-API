package email

import (
	"database/sql"
	"errors"
	"fmt"
	mydb "github.com/wachrusz/Back-End-API/internal/mydatabase"
	"github.com/wachrusz/Back-End-API/internal/myerrors"
	enc "github.com/wachrusz/Back-End-API/pkg/encryption"
	utility "github.com/wachrusz/Back-End-API/pkg/util"

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
	CheckConfirmationCode(email, token, enteredCode string) (CheckResult, error)
	DeleteConfirmationCode(email string, code string) error
	GetConfirmationCode(email string) (string, error)
	ResetPasswordConfirm(token string, code string) (int, int, error)
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
	RemainingAttempts int `json:"remaining_attempts"`
	LockDuration      int `json:"lock_duration"`
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
		return err
	}

	// TODO: email sending

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
func (s *Service) CheckConfirmationCode(email, token, enteredCode string) (CheckResult, error) {
	var result CheckResult
	var expirationTime time.Time
	var attempts int

	err := s.checkToken(token, email)
	if err != nil {
		return result, myerrors.ErrInvalidToken
	}

	locked, err := s.isUserLocked(email)
	if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}
	if locked {
		result.RemainingAttempts = 0
		lockDuration, err := s.getLockDuration(email)
		if err != nil {
			return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
		}
		result.LockDuration = lockDuration
		return result, fmt.Errorf("%w user: try again later", myerrors.ErrLocked)
	}

	err = s.repo.QueryRow(`SELECT attempts FROM confirmation_codes WHERE email = $1 
	ORDER BY expiration_time DESC LIMIT 1`, email).Scan(&attempts)
	if errors.Is(err, sql.ErrNoRows) {
		return result, fmt.Errorf("%w: сan't get attempts", myerrors.ErrInternal)
	} else if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	err = s.repo.QueryRow("SELECT expiration_time FROM confirmation_codes WHERE email = $1 AND code = $2", email, enteredCode).Scan(&expirationTime)
	if errors.Is(err, sql.ErrNoRows) {
		result.RemainingAttempts = maxAttempts - (attempts + 1)

		if result.RemainingAttempts == 0 {
			err := s.lockUser(email)
			if err != nil {
				return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
			}
			lockDuration, err := s.getLockDuration(email)
			if err != nil {
				return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
			}
			result.LockDuration = lockDuration
			return result, myerrors.ErrLocked
		}
		err = s.incrementAttempts(email)
		if err != nil {
			return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
		}
		result.RemainingAttempts = maxAttempts - (attempts + 1)
		return result, myerrors.ErrCode
	} else if err != nil {
		return result, fmt.Errorf("%w: %v", myerrors.ErrInternal, err)
	}

	if time.Now().After(expirationTime) {
		return result, myerrors.ErrExpired
	}

	return result, nil
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

func (s *Service) getLockDuration(email string) (int, error) {
	var lockedUntil int
	query := "SELECT locked_until FROM confirmation_codes WHERE email = $1 ORDER BY locked_until DESC LIMIT 1"
	err := s.repo.QueryRow(query, email).Scan(&lockedUntil)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return lockedUntil, nil
}

// ! DEV/TEST functions
func (s *Service) GetConfirmationCode(email string) (string, error) {
	var code string
	err := s.repo.QueryRow("SELECT code FROM confirmation_codes WHERE email = $1 ORDER BY expiration_time DESC LIMIT 1", email).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *Service) ResetPasswordConfirm(token string, code string) (int, int, error) {
	claims, err := utility.ParseResetToken(token)
	if err != nil {
		return 0, 0, myerrors.ErrInternal
	}

	var registerRequest utility.UserAuthenticationRequest
	registerRequest, err = utility.VerifyResetJWTToken(token)
	if err != nil {
		return 0, 0, myerrors.ErrInvalidToken
	}

	details, err := s.CheckConfirmationCode(registerRequest.Email, token, code)
	if err != nil {
		return details.RemainingAttempts, details.LockDuration, err
	}

	err = s.DeleteConfirmationCode(registerRequest.Email, code)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %v", myerrors.ErrEmailing, err)
	}

	claims["confirmed"] = true
	return 0, 0, nil
}

func (s *Service) ResetPassword(email, password string) error {
	hashedPassword, err := utility.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = s.repo.Exec("UPDATE users SET hashed_password = $1 WHERE email = $2", hashedPassword, email)
	return err
}
