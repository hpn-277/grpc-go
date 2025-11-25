package domain

import "errors"

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailAlreadyExists is returned when trying to create a user with an existing email
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidInput is returned when input validation fails
	ErrInvalidInput = errors.New("invalid input")
)
