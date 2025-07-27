package usecases

import (
	"context"
	"errors"
	"time"

	"ferranrt.com/scheduly-backend/internal/domain"
	"ferranrt.com/scheduly-backend/internal/ports/repositories"
	"ferranrt.com/scheduly-backend/internal/ports/usecases"
	"ferranrt.com/scheduly-backend/internal/services"
	"github.com/google/uuid"
)

type authUseCase struct {
	userRepo        repositories.UserRepository
	sessionRepo     repositories.SessionRepository
	jwtService      services.JWTService
	passwordService services.PasswordService
}

func NewAuthUseCase(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	jwtService services.JWTService,
	passwordService services.PasswordService,
) usecases.AuthUseCase {
	return &authUseCase{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

func (uc *authUseCase) Register(ctx context.Context, registration *domain.UserRegistrationInput, userAgent, ipAddress string) (*domain.AuthResponseDTO, error) {
	// Check if user already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, registration.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := uc.passwordService.HashPassword(registration.Password)
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

	// Generate tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user)
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
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(15 * time.Minute.Seconds()), // 15 minutes
	}, nil
}

func (uc *authUseCase) Login(ctx context.Context, login *domain.UserLoginInput, userAgent, ipAddress string) (*domain.AuthResponseDTO, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = uc.passwordService.VerifyPassword(user.Password, login.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user)
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
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
		ExpiresIn:    int64(15 * time.Minute.Seconds()), // 15 minutes
	}, nil
}

func (uc *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.RefreshTokenResponseDTO, error) {
	// Validate refresh token
	claims, err := uc.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, errors.New("invalid token type")
	}

	// Check if session exists and is active
	session, err := uc.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if !session.IsActive || time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired or inactive")
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	newAccessToken, err := uc.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Update session with new refresh token
	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)
	session.UpdatedAt = time.Now()

	err = uc.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	return &domain.RefreshTokenResponseDTO{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(15 * time.Minute.Seconds()),
	}, nil
}

func (uc *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	// Get session by refresh token
	session, err := uc.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	// Deactivate session
	session.IsActive = false
	session.UpdatedAt = time.Now()

	return uc.sessionRepo.Update(ctx, session)
}

func (uc *authUseCase) LogoutAll(ctx context.Context, userID string) error {
	// Get all active sessions for user
	sessions, err := uc.sessionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Deactivate all sessions
	for _, session := range sessions {
		session.IsActive = false
		session.UpdatedAt = time.Now()
		err = uc.sessionRepo.Update(ctx, session)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*domain.JWTClaims, error) {
	claims, err := uc.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, errors.New("invalid token type")
	}

	// Verify user still exists
	_, err = uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return claims, nil
}

func (uc *authUseCase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserProfileResponseDTO, error) {
	// Get user by ID
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &domain.UserProfileResponseDTO{
		User: *user.ToResponse(),
	}, nil
}
