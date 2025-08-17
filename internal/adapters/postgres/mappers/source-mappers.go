package mappers

import (
	"buke.io/core/internal/adapters/postgres/dbmodels"
	"buke.io/core/internal/domain"
)

type SourceMapper struct {
}

func NewSourceMapper() *SourceMapper {
	return &SourceMapper{}
}

func (m *SourceMapper) DomainToDBModel(source *domain.Source) *dbmodels.Source {
	return &dbmodels.Source{
		ID:                    source.ID,
		CreatedAt:             source.CreatedAt,
		UpdatedAt:             source.UpdatedAt,
		UserID:                source.UserID,
		RefreshToken:          source.RefreshToken,
		UserAgent:             source.UserAgent,
		IPAddress:             source.IPAddress,
		IsActive:              source.IsActive,
		RefreshTokenExpiresAt: source.RefreshTokenExpiresAt,
	}
}

func (m *SourceMapper) DBModelToDomain(source *dbmodels.Source) *domain.Source {
	return &domain.Source{
		ID:                    source.ID,
		UserID:                source.UserID,
		RefreshToken:          source.RefreshToken,
		UserAgent:             source.UserAgent,
		IPAddress:             source.IPAddress,
		IsActive:              source.IsActive,
		RefreshTokenExpiresAt: source.RefreshTokenExpiresAt,
		CreatedAt:             source.CreatedAt,
		UpdatedAt:             source.UpdatedAt,
	}
}
