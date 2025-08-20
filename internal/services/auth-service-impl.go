package services

import (
	"context"
	"errors"
	"time"

	"bifur.app/core/internal/domain"

	"bifur.app/core/internal/exceptions"
	"bifur.app/core/internal/ports"
	"bifur.app/core/internal/utils/password"
	"bifur.app/core/internal/utils/random"
	"bifur.app/core/internal/utils/token"
	"github.com/google/uuid"
)

type AuthServiceImplementation struct {
	userRepo   ports.UserRepository
	sourceRepo ports.SourceRepository
	jwtConfig  domain.JWTConfig
	logger     ports.Logger
}

func NewAuthService(
	userRepo ports.UserRepository,
	sourceRepo ports.SourceRepository,
	jwtConfig domain.JWTConfig,
	logger ports.Logger,
) ports.AuthService {
	return &AuthServiceImplementation{
		userRepo:   userRepo,
		sourceRepo: sourceRepo,
		jwtConfig:  jwtConfig,
		logger:     logger,
	}
}

func generateTokenFromUser(user *domain.User, secret string, duration time.Duration, sourceID string) (string, error) {
	return token.GenerateToken(user.ID, user.Email, []byte(secret), duration, sourceID)
}

func (uc *AuthServiceImplementation) Register(ctx context.Context, input *domain.UserRegisterInput, userAgent, ipAddress string) (*domain.RefreshTokenPayload, error) {
	// Check if user already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		uc.logger.Error(ctx, err)
		return nil, err
	}
	if exists {
		existsError := errors.New("user already exists")
		uc.logger.ErrorWithVar(ctx, existsError, map[string]interface{}{
			"email": input.Email,
		})
		return nil, existsError
	}

	// Hash password
	hashedPassword, err := password.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		uc.logger.Error(ctx, err)
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

	return &domain.RefreshTokenPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(uc.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (uc *AuthServiceImplementation) Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*domain.AuthenticationPayload, error) {
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

	userPayload := &domain.AuthenticationUserPayload{
		ID: user.ID,
	}

	return &domain.AuthenticationPayload{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userPayload,
		ExpiresIn:    int64(uc.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (uc *AuthServiceImplementation) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenPayload, error) {
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
		return nil, exceptions.ErrAuthSourceExpiredOrInvalid
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

	return &domain.RefreshTokenPayload{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(uc.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (uc *AuthServiceImplementation) Logout(ctx context.Context, refreshToken string) error {
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

func (uc *AuthServiceImplementation) LogoutAll(ctx context.Context, userID string) error {
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

func (uc *AuthServiceImplementation) ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error) {
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

func (uc *AuthServiceImplementation) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	// Get user by ID
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
