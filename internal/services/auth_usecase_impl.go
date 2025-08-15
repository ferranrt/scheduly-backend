package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/ports"
	"scheduly.io/core/internal/utils/password"
	"scheduly.io/core/internal/utils/token"
)

type authUseCase struct {
	userRepo    ports.UserRepository
	sessionRepo ports.SessionRepository
	jwtConfig   domain.JWTConfig
}

func NewAuthService(
	userRepo ports.UserRepository,
	sessionRepo ports.SessionRepository,
	jwtConfig domain.JWTConfig,
) ports.AuthService {
	return &authUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtConfig:   jwtConfig,
	}
}

func generateTokenFromUser(user *domain.User, secret string, duration time.Duration) (string, error) {
	return token.GenerateToken(user.ID, user.Email, []byte(secret), duration)
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

	err = useCase.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := generateTokenFromUser(user, useCase.jwtConfig.AtkSecret, useCase.jwtConfig.Expiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateTokenFromUser(user, useCase.jwtConfig.RtkSecret, useCase.jwtConfig.Expiry)
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
		ExpiresAt:    time.Now().Add(useCase.jwtConfig.Expiry),
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
		ExpiresIn:    int64(useCase.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (useCase *authUseCase) Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*domain.AuthResponseDTO, error) {
	// Get user by email
	user, err := useCase.userRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = password.VerifyPassword(user.Password, login.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := generateTokenFromUser(user, useCase.jwtConfig.AtkSecret, useCase.jwtConfig.Expiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateTokenFromUser(user, useCase.jwtConfig.RtkSecret, useCase.jwtConfig.Expiry)
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
		ExpiresAt:    time.Now().Add(useCase.jwtConfig.Expiry),
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
		ExpiresIn:    int64(useCase.jwtConfig.Expiry.Seconds()),
	}, nil
}

func (useCase *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponseDTO, error) {
	// Validate refresh token
	claims, err := token.ValidateToken(refreshToken, []byte(useCase.jwtConfig.AtkSecret))
	if err != nil {
		return nil, errors.New("invalid refresh token")
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
	newAccessToken, err := generateTokenFromUser(user, useCase.jwtConfig.AtkSecret, useCase.jwtConfig.Expiry)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateTokenFromUser(user, useCase.jwtConfig.RtkSecret, useCase.jwtConfig.Expiry)
	if err != nil {
		return nil, err
	}

	// Update session with new refresh token
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(useCase.jwtConfig.Expiry)
	session.UpdatedAt = time.Now()

	err = useCase.sessionRepo.Update(ctx, session)
	if err != nil {

		return nil, err
	}

	return &domain.RefreshTokenResponseDTO{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(useCase.jwtConfig.Expiry.Seconds()),
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

func (useCase *authUseCase) ValidateToken(ctx context.Context, tokenString string) (*domain.JWTClaims, error) {
	claims, err := token.ValidateToken(tokenString, []byte(useCase.jwtConfig.AtkSecret))
	if err != nil {
		return nil, err
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
