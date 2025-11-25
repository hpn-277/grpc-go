package domain

import (
	"github.com/google/uuid"
)

// ID is a generic UUID-based identifier
type ID struct {
	value uuid.UUID
}

// NewID creates a new random ID
func NewID() ID {
	return ID{value: uuid.New()}
}

// ParseID parses a string into an ID
func ParseID(s string) (ID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return ID{}, err
	}
	return ID{value: parsed}, nil
}

// String returns the string representation of the ID
func (id ID) String() string {
	return id.value.String()
}

// Equals checks if two IDs are equal
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}

// IsZero checks if the ID is zero (empty)
func (id ID) IsZero() bool {
	return id.value == uuid.Nil
}
