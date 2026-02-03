package utils

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// ValidateEmail validates an email address format
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// NormalizeEmail normalizes an email address (lowercase, trimmed)
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ValidatePassword validates password requirements
func ValidatePassword(password string) bool {
	// Minimum 8 characters
	return len(password) >= 8
}
