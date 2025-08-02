package dbmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"scheduly.io/core/internal/domain"
)

type Session struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID       uuid.UUID `gorm:"not null;index"`
	RefreshToken string    `gorm:"not null;unique"`
	UserAgent    string    `gorm:"not null"`
	IPAddress    string    `gorm:"not null"`
	IsActive     bool      `gorm:"default:true"`
	ExpiresAt    time.Time `gorm:"not null"`
}

func (s *Session) TableName() string {
	return "sessions"
}

func (s *Session) ToDomain() *domain.Session {
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
