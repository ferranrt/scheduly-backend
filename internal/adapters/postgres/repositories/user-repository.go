package repositories

import (
	"context"
	"errors"

	"scheduly.io/core/internal/adapters/postgres/dbmodels"
	"scheduly.io/core/internal/adapters/postgres/mappers"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/exceptions"
	"scheduly.io/core/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	database   *gorm.DB
	userMapper *mappers.UserMapper
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		database:   db,
		userMapper: mappers.NewUserMapper(),
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	dbUser := r.userMapper.ToDbModel(user)
	result := r.database.WithContext(ctx).Create(dbUser)
	if result.Error != nil {
		return result.Error
	}

	// Update the domain user with the generated ID
	user.ID = dbUser.ID
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var dbUser dbmodels.User
	result := r.database.WithContext(ctx).Where("id = ?", id).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, result.Error
	}

	return r.userMapper.ToDomain(&dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser dbmodels.User
	result := r.database.WithContext(ctx).Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, result.Error
	}

	return r.userMapper.ToDomain(&dbUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	dbUser := r.userMapper.ToDbModel(user)
	result := r.database.WithContext(ctx).Save(dbUser)
	return result.Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.database.WithContext(ctx).Delete(&dbmodels.User{}, "id = ?", id)
	return result.Error
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.database.WithContext(ctx).Model(&dbmodels.User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}
