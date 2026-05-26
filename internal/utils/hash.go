package utils

import (
	
	"encoding/hex"
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(
	password string,
) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func HashSHA256(value string) string {

	hash := sha256.Sum256([]byte(value))

	return hex.EncodeToString(hash[:])
}