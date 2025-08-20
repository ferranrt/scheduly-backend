package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	SourceID string    `json:"source_id"`
	jwt.RegisteredClaims
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type AuthenticationUserPayload struct {
	ID uuid.UUID `json:"id"`
}

type AuthenticationPayload struct {
	AccessToken  string                     `json:"access_token"`
	RefreshToken string                     `json:"refresh_token"`
	User         *AuthenticationUserPayload `json:"user"`
	ExpiresIn    int64                      `json:"expires_in"`
}
