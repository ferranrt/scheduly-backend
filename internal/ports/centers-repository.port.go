package ports

import (
	"context"

	"bifur.app/core/internal/domain"
	"github.com/google/uuid"
)

type CentersRepository interface {
	GetAll(ctx context.Context, userID uuid.UUID) ([]*domain.Center, error)
}
