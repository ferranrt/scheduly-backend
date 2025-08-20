package repositories

import (
	"context"
	"errors"

	"bifur.app/core/internal/adapters/postgres/dbmodels"
	"bifur.app/core/internal/adapters/postgres/mappers"
	"bifur.app/core/internal/domain"
	"bifur.app/core/internal/exceptions"
	"bifur.app/core/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PGUserRepository struct {
	db     *gorm.DB
	mapper *mappers.UserMapper
	logger ports.Logger
}

func NewUserRepository(db *gorm.DB, logger ports.Logger) ports.UserRepository {
	return &PGUserRepository{
		db:     db,
		mapper: mappers.NewUserMapper(),
		logger: logger,
	}
}

func (repo *PGUserRepository) Create(ctx context.Context, user *domain.User) error {
	dbUser := repo.mapper.ToDbModel(user)
	result := repo.db.WithContext(ctx).Create(dbUser)
	if result.Error != nil {
		return result.Error
	}

	// Update the domain user with the generated ID
	user.ID = dbUser.ID
	return nil
}

func (repo *PGUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var dbUser dbmodels.User
	result := repo.db.WithContext(ctx).Where("id = ?", id).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, result.Error
	}

	return repo.mapper.ToDomain(&dbUser), nil
}

func (repo *PGUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser dbmodels.User
	result := repo.db.WithContext(ctx).Where("email = ?", email).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, result.Error
	}

	return repo.mapper.ToDomain(&dbUser), nil
}

func (repo *PGUserRepository) Update(ctx context.Context, user *domain.User) error {
	dbUser := repo.mapper.ToDbModel(user)
	result := repo.db.WithContext(ctx).Save(dbUser)
	return result.Error
}

func (repo *PGUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := repo.db.WithContext(ctx).Delete(&dbmodels.User{}, "id = ?", id)
	return result.Error
}

func (repo *PGUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := repo.db.WithContext(ctx).Model(&dbmodels.User{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}
