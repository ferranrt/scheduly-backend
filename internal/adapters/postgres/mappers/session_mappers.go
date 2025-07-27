package mappers

import (
	"ferranrt.com/scheduly-backend/internal/adapters/postgres/dbmodels"
	"ferranrt.com/scheduly-backend/internal/domain"
)

type SessionMapper struct {
}

func NewSessionMapper() *SessionMapper {
	return &SessionMapper{}
}

func (m *SessionMapper) DomainToDBModel(session *domain.Session) *dbmodels.Session {
	return &dbmodels.Session{
		ID:           session.ID,
		UserID:       session.UserID,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		IPAddress:    session.IPAddress,
		IsActive:     session.IsActive,
		ExpiresAt:    session.ExpiresAt,
	}
}

func (m *SessionMapper) DBModelToDomain(s *dbmodels.Session) *domain.Session {
	return &domain.Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		IPAddress:    s.IPAddress,
		IsActive:     s.IsActive,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
