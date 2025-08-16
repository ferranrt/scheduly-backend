package domain

import (
	"time"

	"github.com/google/uuid"
)

type Source struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	UserAgent             string    `json:"user_agent"`
	IPAddress             string    `json:"ip_address"`
	IsActive              bool      `json:"is_active"`
	City                  string    `json:"city"`
	CountryCode           string    `json:"country_code"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
