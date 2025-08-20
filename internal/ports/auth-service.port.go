package ports

import (
	"context"

	"bifur.app/core/internal/domain"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, registration *domain.UserRegisterInput, userAgent, ipAddress string) (*domain.RefreshTokenPayload, error)
	Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*domain.AuthenticationPayload, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenPayload, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, token string) (*domain.JWTClaims, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}
