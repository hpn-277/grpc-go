package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nguyenphuoc/super-salary-sacrifice/internal/features/user-management/domain"
	sharedDomain "github.com/nguyenphuoc/super-salary-sacrifice/internal/shared/domain"
)

// UserModel is the GORM model for the users table
type UserModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	FirstName string    `gorm:"type:varchar(100);not null"`
	LastName  string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName specifies the table name for GORM
func (UserModel) TableName() string {
	return "users"
}

// GormUserRepository implements domain.UserRepository using GORM
// This is an ADAPTER in hexagonal architecture
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Save persists a user (create or update)
func (r *GormUserRepository) Save(ctx context.Context, user *domain.User) error {
	model := r.toModel(user)
	return r.db.WithContext(ctx).Save(model).Error
}

// FindByID retrieves a user by ID
func (r *GormUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&model)
}

// FindByEmail retrieves a user by email
func (r *GormUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.toDomain(&model)
}

// List retrieves users with pagination
func (r *GormUserRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	var models []UserModel

	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		user, err := r.toDomain(&model)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *GormUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&UserModel{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// toModel converts domain User to GORM UserModel
func (r *GormUserRepository) toModel(user *domain.User) *UserModel {
	id, _ := uuid.Parse(user.ID().String())

	return &UserModel{
		ID:        id,
		Email:     user.Email().String(),
		FirstName: user.FirstName(),
		LastName:  user.LastName(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

// toDomain converts GORM UserModel to domain User
func (r *GormUserRepository) toDomain(model *UserModel) (*domain.User, error) {
	id, err := sharedDomain.ParseID(model.ID.String())
	if err != nil {
		return nil, err
	}

	email, err := sharedDomain.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	return domain.Reconstitute(
		id,
		email,
		model.FirstName,
		model.LastName,
		model.CreatedAt,
		model.UpdatedAt,
	), nil
}
