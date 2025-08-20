package repositories

import (
	"context"

	"bifur.app/core/internal/adapters/postgres/dbmodels"
	"bifur.app/core/internal/adapters/postgres/mappers"
	"bifur.app/core/internal/domain"
	"bifur.app/core/internal/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PGCenterRepository struct {
	db     *gorm.DB
	mapper *mappers.CenterMapper
	logger ports.Logger
}

func NewPgCenterRepository(db *gorm.DB, logger ports.Logger) ports.CentersRepository {
	return &PGCenterRepository{
		db:     db,
		mapper: mappers.NewCenterMapper(),
		logger: logger,
	}
}

func (repo *PGCenterRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]*domain.Center, error) {
	dbCenters := []dbmodels.Center{}
	result := repo.db.WithContext(ctx).Where("owner_id = ?", userID).Find(&dbCenters)
	if result.Error != nil {
		return nil, result.Error
	}

	centers := make([]*domain.Center, len(dbCenters))
	for i, dbCenter := range dbCenters {
		centers[i] = repo.mapper.ToDomain(&dbCenter)
	}
	return centers, nil
}
