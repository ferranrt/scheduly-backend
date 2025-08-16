package repositories

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"scheduly.io/core/internal/adapters/postgres/dbmodels"
	"scheduly.io/core/internal/adapters/postgres/mappers"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/exceptions"
	"scheduly.io/core/internal/ports"
)

type sourceRepository struct {
	database *gorm.DB
	mapper   *mappers.SourceMapper
}

func NewSourceRepository(db *gorm.DB) ports.SourceRepository {
	return &sourceRepository{
		database: db,
		mapper:   mappers.NewSourceMapper(),
	}
}

func (repo *sourceRepository) Create(ctx context.Context, source *domain.Source) error {
	dbSource := repo.mapper.DomainToDBModel(source)
	result := repo.database.WithContext(ctx).Create(dbSource)
	if result.Error != nil {
		return result.Error
	}

	// Update the domain source with the generated ID
	source.ID = dbSource.ID
	return nil
}

func (repo *sourceRepository) GetByID(ctx context.Context, id string) (*domain.Source, error) {
	var dbSource dbmodels.Source
	result := repo.database.WithContext(ctx).Where("id = ?", id).First(&dbSource)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrAuthSourceNotFound
		}
		return nil, result.Error
	}
	return repo.mapper.DBModelToDomain(&dbSource), nil
}

func (repo *sourceRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Source, error) {
	var dbSources []dbmodels.Source
	result := repo.database.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Find(&dbSources)
	if result.Error != nil {
		return nil, result.Error
	}

	sources := make([]*domain.Source, len(dbSources))
	for i, dbSource := range dbSources {
		sources[i] = repo.mapper.DBModelToDomain(&dbSource)
	}

	return sources, nil
}

func (repo *sourceRepository) Update(ctx context.Context, source *domain.Source) error {
	dbSSource := repo.mapper.DomainToDBModel(source)
	result := repo.database.WithContext(ctx).Save(dbSSource)
	return result.Error
}

func (repo *sourceRepository) Delete(ctx context.Context, id string) error {
	result := repo.database.WithContext(ctx).Delete(&dbmodels.Source{}, "id = ?", id)
	return result.Error
}

func (repo *sourceRepository) DeleteByUserID(ctx context.Context, userID string) error {
	result := repo.database.WithContext(ctx).Delete(&dbmodels.Source{}, "user_id = ?", userID)
	return result.Error
}

func (repo *sourceRepository) DeleteExpired(ctx context.Context) error {
	result := repo.database.WithContext(ctx).Delete(&dbmodels.Source{}, "expires_at < ?", time.Now())
	return result.Error
}
