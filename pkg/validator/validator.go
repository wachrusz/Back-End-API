package validator

import (
	"regexp"
)

func IsValidEmail(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func IsValidPassword(password string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9_]{8,30}$`)
	return regex.MatchString(password)
}
