package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidMoney  = errors.New("invalid money amount")
	ErrNegativeMoney = errors.New("money amount cannot be negative")
)

// Money represents a monetary value in cents to avoid floating-point precision issues
// Example: $10.50 is stored as 1050 cents
type Money struct {
	cents int64 // Amount in cents
}

// NewMoney creates a new Money value object from dollars
func NewMoney(dollars float64) (Money, error) {
	if dollars < 0 {
		return Money{}, ErrNegativeMoney
	}

	cents := int64(dollars * 100)
	return Money{cents: cents}, nil
}

// NewMoneyFromCents creates a new Money value object from cents
func NewMoneyFromCents(cents int64) (Money, error) {
	if cents < 0 {
		return Money{}, ErrNegativeMoney
	}

	return Money{cents: cents}, nil
}

// Zero returns a Money value of zero
func Zero() Money {
	return Money{cents: 0}
}

// Cents returns the amount in cents
func (m Money) Cents() int64 {
	return m.cents
}

// Dollars returns the amount in dollars
func (m Money) Dollars() float64 {
	return float64(m.cents) / 100.0
}

// Add adds two Money values
func (m Money) Add(other Money) Money {
	return Money{cents: m.cents + other.cents}
}

// Subtract subtracts two Money values
func (m Money) Subtract(other Money) (Money, error) {
	result := m.cents - other.cents
	if result < 0 {
		return Money{}, ErrNegativeMoney
	}
	return Money{cents: result}, nil
}

// Multiply multiplies Money by a factor
func (m Money) Multiply(factor float64) Money {
	return Money{cents: int64(float64(m.cents) * factor)}
}

// IsZero checks if the amount is zero
func (m Money) IsZero() bool {
	return m.cents == 0
}

// IsGreaterThan checks if this amount is greater than another
func (m Money) IsGreaterThan(other Money) bool {
	return m.cents > other.cents
}

// IsLessThan checks if this amount is less than another
func (m Money) IsLessThan(other Money) bool {
	return m.cents < other.cents
}

// Equals checks if two Money values are equal
func (m Money) Equals(other Money) bool {
	return m.cents == other.cents
}

// String returns a formatted string representation
func (m Money) String() string {
	return fmt.Sprintf("$%.2f", m.Dollars())
}
