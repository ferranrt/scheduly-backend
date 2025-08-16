package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/dtos"
	"scheduly.io/core/internal/exceptions"
	"scheduly.io/core/internal/ports"
	"scheduly.io/core/internal/utils/password"
	"scheduly.io/core/internal/utils/random"
	"scheduly.io/core/internal/utils/token"
)

type authServiceImpl struct {
	userRepo   ports.UserRepository
	sourceRepo ports.SourceRepository
	jwtConfig  domain.JWTConfig
}

func NewAuthService(
	userRepo ports.UserRepository,
	sourceRepo ports.SourceRepository,
	jwtConfig domain.JWTConfig,
) ports.AuthService {
	return &authServiceImpl{
		userRepo:   userRepo,
		sourceRepo: sourceRepo,
		jwtConfig:  jwtConfig,
	}
}

func generateTokenFromUser(user *domain.User, secret string, duration time.Duration, sourceID string) (string, error) {
	return token.GenerateToken(user.ID, user.Email, []byte(secret), duration, sourceID)
}

func (uc *authServiceImpl) Register(ctx context.Context, registration *domain.UserRegisterInput, userAgent, ipAddress string) (*dtos.AuthResponseDTO, error) {
	// Check if user already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, registration.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := password.HashPassword(registration.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:     registration.Email,
		Password:  hashedPassword,
		FirstName: registration.FirstName,
		LastName:  registration.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	refreshToken := random.GenerateRandomString(128)
	newSource := &domain.Source{
		UserID:                user.ID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: time.Now().Add(uc.jwtConfig.Expiry),
		UserAgent:             userAgent,

		IPAddress: ipAddress,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.sourceRepo.Create(ctx, newSource)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := generateTokenFromUser(user, uc.jwtConfig.AtkSecret, uc.jwtConfig.Expiry, newSource.ID.String())
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(uc.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (uc *authServiceImpl) Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*dtos.AuthResponseDTO, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, exceptions.ErrAuthInvalidCredentials
	}

	// Verify password
	err = password.VerifyPassword(user.Password, login.Password)
	if err != nil {
		return nil, exceptions.ErrAuthInvalidCredentials
	}

	refreshToken := random.GenerateRandomString(128)

	newSource := &domain.Source{
		UserID:                user.ID,
		RefreshToken:          refreshToken,
		UserAgent:             userAgent,
		IPAddress:             ipAddress,
		IsActive:              true,
		RefreshTokenExpiresAt: time.Now().Add(uc.jwtConfig.Expiry),
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	err = uc.sourceRepo.Create(ctx, newSource)
	if err != nil {
		return nil, err
	}
	// Generate tokens
	accessToken, err := generateTokenFromUser(user, uc.jwtConfig.AtkSecret, uc.jwtConfig.Expiry, newSource.ID.String())
	if err != nil {
		return nil, err
	}

	return &dtos.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(uc.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (uc *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponseDTO, error) {
	// Validate refresh token
	claims, err := token.ValidateToken(refreshToken, []byte(uc.jwtConfig.AtkSecret))
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if source exists and is active
	source, err := uc.sourceRepo.GetByID(ctx, claims.SourceID)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if !source.IsActive || time.Now().After(source.RefreshTokenExpiresAt) {
		return nil, exceptions.ErrSourceExpiredOrInvalid
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	newAccessToken, err := generateTokenFromUser(user, uc.jwtConfig.AtkSecret, uc.jwtConfig.Expiry, source.ID.String())
	if err != nil {
		return nil, err
	}

	newRefreshToken := random.GenerateRandomString(128)

	source.RefreshToken = newRefreshToken
	source.RefreshTokenExpiresAt = time.Now().Add(uc.jwtConfig.Expiry)
	source.UpdatedAt = time.Now()

	err = uc.sourceRepo.Update(ctx, source)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshTokenResponseDTO{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(uc.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (uc *authServiceImpl) Logout(ctx context.Context, refreshToken string) error {
	// Get source by refresh token
	source, err := uc.sourceRepo.GetByID(ctx, refreshToken)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	// Deactivate source
	source.IsActive = false
	source.UpdatedAt = time.Now()

	return uc.sourceRepo.Update(ctx, source)
}

func (uc *authServiceImpl) LogoutAll(ctx context.Context, userID string) error {
	// Get all active sources for user
	sources, err := uc.sourceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Deactivate all sources
	for _, source := range sources {
		source.IsActive = false
		source.UpdatedAt = time.Now()
		err = uc.sourceRepo.Update(ctx, source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *authServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error) {
	claims, err := token.ValidateToken(tokenString, []byte(uc.jwtConfig.AtkSecret))
	if err != nil {
		return nil, err
	}

	// Verify user still exists
	_, err = uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return claims, nil
}

func (uc *authServiceImpl) GetProfile(ctx context.Context, userID uuid.UUID) (*dtos.UserProfileResponseDTO, error) {
	// Get user by ID
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dtos.UserProfileResponseDTO{
		User: *user.ToResponse(),
	}, nil
}
