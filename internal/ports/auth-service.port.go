package ports

import (
	"context"

	"buke.io/core/internal/domain"
	"buke.io/core/internal/dtos"
	"github.com/google/uuid"
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
