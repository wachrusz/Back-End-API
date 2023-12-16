package email

import (
	"time"

	enc "backEndAPI/_encryption"
	logger "backEndAPI/_logger"
	mydb "backEndAPI/_mydatabase"
	utility "backEndAPI/_utility"
	//"github.com/go-gomail/gomail"
)

var (
	corpEmail         string = "ulohirapshit@gmail.com"
	corpEmailPassword string = "Ihatetechies322"
	webURL            string = ""
)

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

func CheckConfirmationCode(email, enteredCode string) bool {
	var expirationTime time.Time
	err := mydb.GlobalDB.QueryRow("SELECT expiration_time FROM confirmation_codes WHERE email = $1 AND code = $2", email, enteredCode).Scan(&expirationTime)
	if err != nil {
		logger.ErrorLogger.Printf("Confirmation code not found for UserID: %v", email)
		return false
	}

	return time.Now().Before(expirationTime)
}

func DeleteConfimationCode(email string) {
	err := mydb.GlobalDB.QueryRow("DELETE FROM confirmation_codes WHERE email = $1", email)
	if err != nil {
		logger.ErrorLogger.Printf("Error deleting code for Email: %v", email)
		return
	}
}

func ConfirmEmail(email string) error {
	DeleteConfimationCode(email)
	return nil
}
