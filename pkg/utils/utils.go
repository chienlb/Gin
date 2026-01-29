package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func VerifyPassword(hashedPassword, password string) bool {
	return HashPassword(password) == hashedPassword
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}
