package email

import (
	"database/sql"
	json "encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	enc "main/packages/_encryption"

	jsonresponse "main/packages/_json_response"

	logger "main/packages/_logger"
	mydb "main/packages/_mydatabase"
	utility "main/packages/_utility"
	//"github.com/go-gomail/gomail"
)

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

func SendEmail(to, subject, body string) error { /*
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

func SendConfirmationEmail(email, token string) error {
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
	SaveConfirmationCode(email, confirmationCode, token)

	return nil
}

func SendConfirmationEmailTestHandler(email, token string, w http.ResponseWriter, r *http.Request) error {
	confirmationCode, err := utility.GenerateConfirmationCode()
	if err != nil {
		logger.ErrorLogger.Printf("Error in generating confirmation code for Email: %v", email)
		return err
	}

	SaveConfirmationCode(email, confirmationCode, token)

	response := map[string]interface{}{
		"message":     "Successfully sent confirmation code.",
		"code":        confirmationCode,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	return nil
}

func GetConfirmationCodeTestHandler(w http.ResponseWriter, r *http.Request) {
	type jsonEmail struct {
		Email string `json:"email"`
	}

	//! ELDER VER
	/*
		var email_struct jsonEmail

			err := json.NewDecoder(r.Body).Decode(&email_struct)
			if err != nil {
				jsonresponse.SendErrorResponse(w, errors.New("Invalid request payload: "+err.Error()), http.StatusBadRequest)
				return
			}
			email := email_struct.Email
	*/

	email := r.URL.Query().Get("email")
	if email == "" {
		jsonresponse.SendErrorResponse(w, errors.New("Incorrect email"), http.StatusBadRequest)
		return
	}

	code := getConfirmationCode(email)
	if code == "" {
		jsonresponse.SendErrorResponse(w, errors.New("Email not found."), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"message":     "Successfully sent confirmation code.",
		"code":        code,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DecryptToken(encryptedToken string) (string, error) {
	return enc.DecryptToken(encryptedToken)
}

func SaveConfirmationCode(email, confirmationCode, token string) error {
	encryptedToken, errEnc := enc.EncryptToken(token)
	if errEnc != nil {
		return errEnc
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	_, err := mydb.GlobalDB.Exec("INSERT INTO confirmation_codes (email, code, expiration_time, token) VALUES ($1, $2, $3, $4)", email, confirmationCode, expirationTime, encryptedToken)
	return err
}

// ! ДОДЕЛАТЬ
func CheckConfirmationCode(email, token, enteredCode string) CheckResult {
	var result CheckResult
	var expirationTime time.Time
	var attempts int
	result.StatusCode = http.StatusOK
	result.Err = "nil"

	err := checkToken(token, email)
	if err != nil {
		fmt.Println(err)
		result.Err = errors.New("Invalid token for the last code").Error()
		result.StatusCode = http.StatusUnauthorized
		return result
	}

	locked, err := isUserLocked(email)
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
		lockDuration, err := getLockDuration(email)
		if err != nil {
			result.Err = errors.New("Server error").Error()
		}
		result.LockDuration = lockDuration
		result.StatusCode = http.StatusUnauthorized
		return result
	}

	err = mydb.GlobalDB.QueryRow(`SELECT attempts FROM confirmation_codes WHERE email = $1 
	ORDER BY expiration_time DESC LIMIT 1`, email).Scan(&attempts)
	if err == sql.ErrNoRows {
		result.Err = errors.New("Can't get attempts").Error()
		result.StatusCode = http.StatusUnauthorized
		return result
	} else if err != nil {
		result.Err = err.Error()
		result.StatusCode = http.StatusUnauthorized
		return result
	}

	err = mydb.GlobalDB.QueryRow("SELECT expiration_time FROM confirmation_codes WHERE email = $1 AND code = $2", email, enteredCode).Scan(&expirationTime)
	if err == sql.ErrNoRows {
		result.Err = errors.New("Invalid confirmation code.").Error()
		result.RemainingAttempts = maxAttempts - (attempts + 1)

		if maxAttempts-(attempts+1) == 0 {
			err := lockUser(email)
			if err != nil {
				result.Err = err.Error()
				return result
			}

			result.Err = errors.New("User is locked. Try again later.").Error()
			lockDuration, err := getLockDuration(email)
			if err != nil {
				result.Err = errors.New("Server error").Error()
			}
			result.LockDuration = lockDuration
			result.StatusCode = http.StatusUnauthorized
			return result
		}

		err = incrementAttempts(email)
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

func checkToken(token, email string) error {
	var encToken string
	err := mydb.GlobalDB.QueryRow(`SELECT token FROM confirmation_codes WHERE email = $1
	ORDER BY expiration_time DESC LIMIT 1`, email).Scan(&encToken)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	decToken, err := DecryptToken(encToken)
	if err != nil {
		return err
	}

	if token != decToken {
		return errors.New("Invalid token for code")
	}
	return nil
}

func DeleteConfirmationCode(email string, code string) {
	err := mydb.GlobalDB.QueryRow("DELETE FROM confirmation_codes WHERE email = $1 AND code = $2", email, code)
	if err != nil {
		logger.ErrorLogger.Printf("Error deleting code for Email: %v", email)
		return
	}
}

func ConfirmEmail(email, code string) error {
	DeleteConfirmationCode(email, code)
	return nil
}

func incrementAttempts(email string) error {
	query := `UPDATE confirmation_codes SET attempts = attempts + 1 WHERE email = $1 AND expiration_time = (
		SELECT MAX(expiration_time) 
		FROM confirmation_codes 
		WHERE email = $1
	);`
	_, err := mydb.GlobalDB.Exec(query, email)
	return err
}

func lockUser(email string) error {
	query := `UPDATE confirmation_codes SET attempts = 0, 
	locked_until = NOW() + (5 * interval '1 minute') WHERE email = $1 
		AND expiration_time = (
		SELECT MAX(expiration_time) 
		FROM confirmation_codes 
		WHERE email = $1
	);`
	_, err := mydb.GlobalDB.Exec(query, email)
	return err
}

func isUserLocked(email string) (bool, error) {
	query := "SELECT locked_until FROM confirmation_codes WHERE email = $1 ORDER BY locked_until DESC LIMIT 1"
	var lockedUntil time.Time
	err := mydb.GlobalDB.QueryRow(query, email).Scan(&lockedUntil)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return time.Now().Before(lockedUntil), nil
}

func getLockDuration(email string) (time.Duration, error) {
	var lockedUntil time.Time

	query := "SELECT locked_until FROM confirmation_codes WHERE email = $1 ORDER BY locked_until DESC LIMIT 1"
	err := mydb.GlobalDB.QueryRow(query, email).Scan(&lockedUntil)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return time.Until(lockedUntil), nil
}

// ! DEV/TEST functions
func getConfirmationCode(email string) string {
	var code string
	err := mydb.GlobalDB.QueryRow("SELECT code FROM confirmation_codes WHERE email = $1 ORDER BY expiration_time DESC LIMIT 1", email).Scan(&code)
	if err != nil {
		logger.ErrorLogger.Printf("Confirmation code not found for Email: %v", email)
		return ""
	}
	return code
}
