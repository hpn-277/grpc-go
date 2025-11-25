package domain

import "context"

// UserRepository defines the port (interface) for user persistence
// This is a PORT in hexagonal architecture - the domain defines what it needs
type UserRepository interface {
	// Save persists a user (create or update)
	Save(ctx context.Context, user *User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id string) (*User, error)

	// FindByEmail retrieves a user by email
	FindByEmail(ctx context.Context, email string) (*User, error)

	// List retrieves users with pagination
	List(ctx context.Context, offset, limit int) ([]*User, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
