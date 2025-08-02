package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/ports/repositories"
	"scheduly.io/core/internal/ports/usecases"
	"scheduly.io/core/internal/services"
)

type authUseCase struct {
	userRepo        repositories.UserRepository
	sessionRepo     repositories.SessionRepository
	jwtService      services.JWTService
	passwordService services.PasswordService
	jwtConfig       domain.JWTConfig
}

func NewAuthUseCase(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	jwtService services.JWTService,
	passwordService services.PasswordService,
	jwtConfig domain.JWTConfig,
) usecases.AuthUseCase {
	return &authUseCase{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
		jwtConfig:       jwtConfig,
	}
}

func (useCase *authUseCase) Register(ctx context.Context, registration *domain.UserRegistrationInput, userAgent, ipAddress string) (*domain.AuthResponseDTO, error) {
	// Check if user already exists
	exists, err := useCase.userRepo.ExistsByEmail(ctx, registration.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := useCase.passwordService.HashPassword(registration.Password)
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

	err = useCase.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := useCase.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := useCase.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(useCase.jwtConfig.RefreshTokenExpiry),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = useCase.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(useCase.jwtConfig.AccessTokenExpiry.Seconds()),
	}, nil
}

func (useCase *authUseCase) Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*domain.AuthResponseDTO, error) {
	// Get user by email
	user, err := useCase.userRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = useCase.passwordService.VerifyPassword(user.Password, login.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := useCase.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := useCase.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &domain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(useCase.jwtConfig.RefreshTokenExpiry),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = useCase.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(useCase.jwtConfig.AccessTokenExpiry.Seconds()),
	}, nil
}

func (useCase *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponseDTO, error) {
	// Validate refresh token
	claims, err := useCase.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Check if session exists and is active
	session, err := useCase.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired or inactive")
	}

	// Get user
	user, err := useCase.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	newAccessToken, err := useCase.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := useCase.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Update session with new refresh token
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(useCase.jwtConfig.RefreshTokenExpiry)
	session.UpdatedAt = time.Now()

	err = useCase.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshTokenResponseDTO{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(useCase.jwtConfig.AccessTokenExpiry.Seconds()),
	}, nil
}

func (useCase *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	// Get session by refresh token
	session, err := useCase.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	// Deactivate session
	session.IsActive = false
	session.UpdatedAt = time.Now()

	return useCase.sessionRepo.Update(ctx, session)
}

func (useCase *authUseCase) LogoutAll(ctx context.Context, userID string) error {
	// Get all active sessions for user
	sessions, err := useCase.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Deactivate all sessions
	for _, session := range sessions {
		session.IsActive = false
		session.UpdatedAt = time.Now()
		err = useCase.sessionRepo.Update(ctx, session)
		if err != nil {
			return err
		}
	}

	return nil
}

func (useCase *authUseCase) ValidateToken(ctx context.Context, token string) (*domain.JWTClaims, error) {
	claims, err := useCase.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, errors.New("invalid token type")
	}

	// Verify user still exists
	_, err = useCase.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return claims, nil
}

func (useCase *authUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfileResponseDTO, error) {
	// Get user by ID
	user, err := useCase.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &domain.UserProfileResponseDTO{
		User: *user.ToResponse(),
	}, nil
}
