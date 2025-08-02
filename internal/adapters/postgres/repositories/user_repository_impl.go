package repositories

import (
	"context"
	"errors"

	"scheduly.io/core/internal/adapters/postgres/dbmodels"
	"scheduly.io/core/internal/adapters/postgres/mappers"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/ports/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	database *gorm.DB
	mapper   *mappers.UserMapper
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{database: db, mapper: mappers.NewUserMapper()}
}

func (repo *userRepository) Create(ctx context.Context, user *domain.User) error {
	dbUser := repo.mapper.DomainToDBModel(user)
	result := repo.database.WithContext(ctx).Create(dbUser)
	if result.Error != nil {
		return result.Error
	}

	// Update the domain user with the generated ID
	user.ID = dbUser.ID
	return nil
}

func (repo *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var dbUser dbmodels.User
	result := repo.database.WithContext(ctx).Where("id = ?", id).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return repo.mapper.DBModelToDomain(&dbUser), nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser dbmodels.User
	result := repo.database.WithContext(ctx).Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	return repo.mapper.DBModelToDomain(&dbUser), nil
}

func (repo *userRepository) Update(ctx context.Context, user *domain.User) error {
	dbUser := repo.mapper.DomainToDBModel(user)
	result := repo.database.WithContext(ctx).Save(dbUser)
	return result.Error
}

func (repo *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := repo.database.WithContext(ctx).Delete(&dbmodels.User{}, "id = ?", id)
	return result.Error
}

func (repo *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := repo.database.WithContext(ctx).Model(&dbmodels.User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}
