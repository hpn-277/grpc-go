package domain

import (
	"time"

	"github.com/nguyenphuoc/super-salary-sacrifice/internal/shared/domain"
)

// User represents a user in the system (aggregate root)
type User struct {
	id        domain.ID
	email     domain.Email
	firstName string
	lastName  string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser creates a new User
func NewUser(email domain.Email, firstName, lastName string) *User {
	now := time.Now().UTC()
	return &User{
		id:        domain.NewID(),
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		createdAt: now,
		updatedAt: now,
	}
}

// Reconstitute creates a User from persistence (used by repository)
func Reconstitute(
	id domain.ID,
	email domain.Email,
	firstName, lastName string,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:        id,
		email:     email,
		firstName: firstName,
		lastName:  lastName,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// Getters (read-only access to encapsulated fields)
func (u *User) ID() domain.ID       { return u.id }
func (u *User) Email() domain.Email { return u.email }
func (u *User) FirstName() string   { return u.firstName }
func (u *User) LastName() string    { return u.lastName }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// Business methods

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.firstName + " " + u.lastName
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(firstName, lastName string) {
	u.firstName = firstName
	u.lastName = lastName
	u.updatedAt = time.Now().UTC()
}

// ChangeEmail changes the user's email address
func (u *User) ChangeEmail(newEmail domain.Email) {
	u.email = newEmail
	u.updatedAt = time.Now().UTC()
}
