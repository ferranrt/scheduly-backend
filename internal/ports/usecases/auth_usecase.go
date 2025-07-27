package usecases

import (
	"context"

	"ferranrt.com/scheduly-backend/internal/domain"
	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(ctx context.Context, registration *domain.UserRegistrationInput, userAgent, ipAddress string) (*domain.AuthResponseDTO, error)
	Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*domain.AuthResponseDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponseDTO, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, token string) (*domain.JWTClaims, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfileResponseDTO, error)
}
