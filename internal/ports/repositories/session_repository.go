package repositories

import (
	"context"

	"scheduly.io/core/internal/domain"
)

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	GetByUserID(ctx context.Context, userID string) ([]*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}
