package repositories

import (
	"context"
	"errors"
	"time"

	"buke.io/core/internal/adapters/postgres/dbmodels"
	"buke.io/core/internal/adapters/postgres/mappers"
	"buke.io/core/internal/domain"
	"buke.io/core/internal/exceptions"
	"buke.io/core/internal/ports"
	"gorm.io/gorm"
)

type PGSourceRepository struct {
	db     *gorm.DB
	mapper *mappers.SourceMapper
	logger ports.Logger
}

func NewSourceRepository(db *gorm.DB, logger ports.Logger) ports.SourceRepository {
	return &PGSourceRepository{
		db:     db,
		mapper: mappers.NewSourceMapper(),
		logger: logger,
	}
}

func (repo *PGSourceRepository) Create(ctx context.Context, source *domain.Source) error {
	dbSource := repo.mapper.DomainToDBModel(source)
	result := repo.db.WithContext(ctx).Create(dbSource)
	if result.Error != nil {
		return result.Error
	}

	// Update the domain source with the generated ID
	source.ID = dbSource.ID
	return nil
}

func (repo *PGSourceRepository) GetByID(ctx context.Context, id string) (*domain.Source, error) {
	var dbSource dbmodels.Source
	result := repo.db.WithContext(ctx).Where("id = ?", id).First(&dbSource)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrAuthSourceNotFound
		}
		return nil, result.Error
	}
	return repo.mapper.DBModelToDomain(&dbSource), nil
}

func (repo *PGSourceRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Source, error) {
	var dbSources []dbmodels.Source
	result := repo.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Find(&dbSources)
	if result.Error != nil {
		return nil, result.Error
	}

	sources := make([]*domain.Source, len(dbSources))
	for i, dbSource := range dbSources {
		sources[i] = repo.mapper.DBModelToDomain(&dbSource)
	}

	return sources, nil
}

func (repo *PGSourceRepository) Update(ctx context.Context, source *domain.Source) error {
	dbSSource := repo.mapper.DomainToDBModel(source)
	result := repo.db.WithContext(ctx).Save(dbSSource)
	return result.Error
}

func (repo *PGSourceRepository) Delete(ctx context.Context, id string) error {
	result := repo.db.WithContext(ctx).Delete(&dbmodels.Source{}, "id = ?", id)
	return result.Error
}

func (repo *PGSourceRepository) DeleteByUserID(ctx context.Context, userID string) error {
	result := repo.db.WithContext(ctx).Delete(&dbmodels.Source{}, "user_id = ?", userID)
	return result.Error
}

func (repo *PGSourceRepository) DeleteExpired(ctx context.Context) error {
	result := repo.db.WithContext(ctx).Delete(&dbmodels.Source{}, "expires_at < ?", time.Now())
	return result.Error
}
