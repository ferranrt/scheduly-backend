package mappers

import (
	"bifur.app/core/internal/adapters/postgres/dbmodels"
	"bifur.app/core/internal/domain"
)

type CenterMapper struct{}

func NewCenterMapper() *CenterMapper {
	return &CenterMapper{}
}

func (m *CenterMapper) ToDbModel(center *domain.Center) *dbmodels.Center {
	return &dbmodels.Center{
		Name:    center.Name,
		ID:      center.ID,
		OwnerID: center.OwnerID,
	}
}

func (m *CenterMapper) ToDomain(center *dbmodels.Center) *domain.Center {
	return &domain.Center{
		ID:      center.ID,
		Name:    center.Name,
		OwnerID: center.OwnerID,
	}
}
