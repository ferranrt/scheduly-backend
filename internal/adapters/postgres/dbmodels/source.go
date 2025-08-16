package dbmodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"scheduly.io/core/internal/domain"
)

type Source struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt `gorm:"index"`
	IPAddress             string         `gorm:"not null"`
	UserID                uuid.UUID      `gorm:"not null;index"`
	RefreshToken          string         `gorm:"not null;unique"`
	RefreshTokenExpiresAt time.Time      `gorm:"not null"`
	UserAgent             string         `gorm:"not null"`
	IsActive              bool           `gorm:"default:true"`
	City                  string         `gorm:"not null"`
	CountryCode           string         `gorm:"not null"`
}

func (s *Source) TableName() string {
	return "sources"
}

func (s *Source) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

func (s *Source) ToDomain() *domain.Source {
	return &domain.Source{
		ID:                    s.ID,
		CreatedAt:             s.CreatedAt,
		UpdatedAt:             s.UpdatedAt,
		UserID:                s.UserID,
		RefreshToken:          s.RefreshToken,
		RefreshTokenExpiresAt: s.RefreshTokenExpiresAt,
		UserAgent:             s.UserAgent,
		IPAddress:             s.IPAddress,
		IsActive:              s.IsActive,
		City:                  s.City,
		CountryCode:           s.CountryCode,
	}
}
