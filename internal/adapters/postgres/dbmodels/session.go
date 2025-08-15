package dbmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"scheduly.io/core/internal/domain"
)

type Session struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	UserID       uuid.UUID      `gorm:"not null;index"`
	RefreshToken string         `gorm:"not null;unique"`
	UserAgent    string         `gorm:"not null"`
	IPAddress    string         `gorm:"not null"`
	IsActive     bool           `gorm:"default:true"`
	ExpiresAt    time.Time      `gorm:"not null"`
}

func (s *Session) TableName() string {
	return "sessions"
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
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
