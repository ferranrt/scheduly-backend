package ports

import (
	"context"

	"github.com/google/uuid"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/dtos"
)

type AuthService interface {
	Register(ctx context.Context, registration *domain.UserRegisterInput, userAgent, ipAddress string) (*dtos.AuthResponseDTO, error)
	Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*dtos.AuthResponseDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponseDTO, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, token string) (*domain.JWTClaims, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*dtos.UserProfileResponseDTO, error)
}
