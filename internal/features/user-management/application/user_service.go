package application

import (
	"context"
	"fmt"

	"github.com/nguyenphuoc/super-salary-sacrifice/internal/features/user-management/domain"
	sharedDomain "github.com/nguyenphuoc/super-salary-sacrifice/internal/shared/domain"
)

// Commands (write operations)

// CreateUserCommand represents the command to create a new user
type CreateUserCommand struct {
	Email     string
	FirstName string
	LastName  string
}

// Queries (read operations)

// GetUserQuery represents a query to get a user by ID
type GetUserQuery struct {
	UserID string
}

// ListUsersQuery represents a query to list users with pagination
type ListUsersQuery struct {
	Offset int
	Limit  int
}

// UserService handles user management use cases
// This is the APPLICATION layer - orchestrates domain logic and repository
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user (use case)
func (s *UserService) CreateUser(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
	// Validate email using domain value object
	email, err := sharedDomain.NewEmail(cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, email.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Validate input
	if cmd.FirstName == "" || cmd.LastName == "" {
		return nil, domain.ErrInvalidInput
	}

	// Create user using domain factory
	user := domain.NewUser(email, cmd.FirstName, cmd.LastName)

	// Persist user
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID (use case)
func (s *UserService) GetUser(ctx context.Context, query GetUserQuery) (*domain.User, error) {
	// Validate ID format
	if query.UserID == "" {
		return nil, domain.ErrInvalidInput
	}

	// Retrieve user from repository
	user, err := s.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers retrieves a paginated list of users (use case)
func (s *UserService) ListUsers(ctx context.Context, query ListUsersQuery) ([]*domain.User, error) {
	// Set default pagination if not provided
	if query.Limit <= 0 {
		query.Limit = 10 // Default page size
	}
	if query.Limit > 100 {
		query.Limit = 100 // Max page size
	}
	if query.Offset < 0 {
		query.Offset = 0
	}

	// Retrieve users from repository
	users, err := s.userRepo.List(ctx, query.Offset, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
