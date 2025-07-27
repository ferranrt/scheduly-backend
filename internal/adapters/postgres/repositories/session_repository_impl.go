package repositories

import (
	"context"
	"errors"
	"time"

	"ferranrt.com/scheduly-backend/internal/adapters/postgres/dbmodels"
	"ferranrt.com/scheduly-backend/internal/adapters/postgres/mappers"
	"ferranrt.com/scheduly-backend/internal/ports/repositories"

	"ferranrt.com/scheduly-backend/internal/domain"
	"gorm.io/gorm"
)

type sessionRepository struct {
	database *gorm.DB
	mapper   *mappers.SessionMapper
}

func NewSessionRepository(db *gorm.DB) repositories.SessionRepository {
	return &sessionRepository{
		database: db,
		mapper:   mappers.NewSessionMapper(),
	}
}

func (repo *sessionRepository) Create(ctx context.Context, session *domain.Session) error {
	dbSession := repo.mapper.DomainToDBModel(session)
	result := repo.database.WithContext(ctx).Create(dbSession)
	if result.Error != nil {
		return result.Error
	}

	// Update the domain session with the generated ID
	session.ID = dbSession.ID
	return nil
}

func (repo *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	var dbSession dbmodels.Session
	result := repo.database.WithContext(ctx).Where("refresh_token = ? AND is_active = ?", refreshToken, true).First(&dbSession)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("session not found")
		}
		return nil, result.Error
	}

	return repo.mapper.DBModelToDomain(&dbSession), nil
}

func (repo *sessionRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Session, error) {
	var dbSessions []dbmodels.Session
	result := repo.database.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Find(&dbSessions)
	if result.Error != nil {
		return nil, result.Error
	}

	sessions := make([]*domain.Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = repo.mapper.DBModelToDomain(&dbSession)
	}

	return sessions, nil
}

func (repo *sessionRepository) Update(ctx context.Context, session *domain.Session) error {
	dbSession := repo.mapper.DomainToDBModel(session)
	result := repo.database.WithContext(ctx).Save(dbSession)
	return result.Error
}

func (repo *sessionRepository) Delete(ctx context.Context, id string) error {
	result := repo.database.WithContext(ctx).Delete(&dbmodels.Session{}, "id = ?", id)
	return result.Error
}

func (repo *sessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	result := repo.database.WithContext(ctx).Delete(&dbmodels.Session{}, "user_id = ?", userID)
	return result.Error
}

func (repo *sessionRepository) DeleteExpired(ctx context.Context) error {
	result := repo.database.WithContext(ctx).Delete(&dbmodels.Session{}, "expires_at < ?", time.Now())
	return result.Error
}
