package validation

import (
	"net/mail"
	"regexp"

	"unicode"

	"github.com/google/uuid"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidIndianPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(
		`^\+91[6-9][0-9]{9}$`,
	)

	return phoneRegex.MatchString(phone)
}

func IsStrongPassword(password string) bool {

	if len(password) < 8 {
		return false
	}

	var (
		hasUpper bool
		hasLower bool
		hasDigit bool
	)

	for _, char := range password {

		switch {
		case unicode.IsUpper(char):
			hasUpper = true

		case unicode.IsLower(char):
			hasLower = true

		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

func IsValidUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}