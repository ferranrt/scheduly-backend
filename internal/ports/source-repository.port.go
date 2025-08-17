package ports

import (
	"context"

	"buke.io/core/internal/domain"
)

type SourceRepository interface {
	Create(ctx context.Context, source *domain.Source) error
	GetByUserID(ctx context.Context, userID string) ([]*domain.Source, error)
	GetByID(ctx context.Context, id string) (*domain.Source, error)
	Update(ctx context.Context, source *domain.Source) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
}
