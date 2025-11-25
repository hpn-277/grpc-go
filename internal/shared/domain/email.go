package domain

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidEmail = errors.New("invalid email format")
	emailRegex      = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// Email is a value object representing a validated email address
type Email struct {
	value string
}

// NewEmail creates a new Email value object with validation
func NewEmail(email string) (Email, error) {
	// Normalize: trim whitespace and convert to lowercase
	normalized := strings.TrimSpace(strings.ToLower(email))

	if normalized == "" {
		return Email{}, ErrInvalidEmail
	}

	if !emailRegex.MatchString(normalized) {
		return Email{}, ErrInvalidEmail
	}

	return Email{value: normalized}, nil
}

// String returns the email as a string
func (e Email) String() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsEmpty checks if the email is empty
func (e Email) IsEmpty() bool {
	return e.value == ""
}
