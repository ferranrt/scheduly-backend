package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"ferranrt.com/scheduly-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations for testing
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Session, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Session), args.Error(1)
}

func (m *MockSessionRepository) Update(ctx context.Context, session *domain.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccessToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(user *domain.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*domain.JWTClaims), args.Error(1)
}

func (m *MockJWTService) GenerateRandomToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func TestAuthUseCase_Register(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)

	authUseCase := NewAuthUseCase(mockUserRepo, mockSessionRepo, mockJWTService, mockPasswordService)

	registration := &domain.UserRegistrationInput{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Expectations
	mockUserRepo.On("ExistsByEmail", mock.Anything, registration.Email).Return(false, nil)
	mockPasswordService.On("HashPassword", registration.Password).Return("hashed_password", nil)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
	mockJWTService.On("GenerateAccessToken", mock.AnythingOfType("*domain.User")).Return("access_token", nil)
	mockJWTService.On("GenerateRefreshToken", mock.AnythingOfType("*domain.User")).Return("refresh_token", nil)
	mockSessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Session")).Return(nil)

	// Execute
	response, err := authUseCase.Register(context.Background(), registration, "test-agent", "127.0.0.1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "access_token", response.AccessToken)
	assert.Equal(t, "refresh_token", response.RefreshToken)
	assert.NotNil(t, response.User)
	assert.Equal(t, registration.Email, response.User.Email)

	// Verify all expectations were met
	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	mockJWTService.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
}

func TestAuthUseCase_Register_UserAlreadyExists(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)

	authUseCase := NewAuthUseCase(mockUserRepo, mockSessionRepo, mockJWTService, mockPasswordService)

	registration := &domain.UserRegistrationInput{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	// Expectations
	mockUserRepo.On("ExistsByEmail", mock.Anything, registration.Email).Return(true, nil)

	// Execute
	response, err := authUseCase.Register(context.Background(), registration, "test-agent", "127.0.0.1")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "user already exists", err.Error())

	// Verify all expectations were met
	mockUserRepo.AssertExpectations(t)
}

func TestAuthUseCase_Login(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)

	authUseCase := NewAuthUseCase(mockUserRepo, mockSessionRepo, mockJWTService, mockPasswordService)

	login := &domain.UserLoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	user := &domain.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  "hashed_password",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Expectations
	mockUserRepo.On("GetByEmail", mock.Anything, login.Email).Return(user, nil)
	mockPasswordService.On("VerifyPassword", user.Password, login.Password).Return(nil)
	mockJWTService.On("GenerateAccessToken", user).Return("access_token", nil)
	mockJWTService.On("GenerateRefreshToken", user).Return("refresh_token", nil)
	mockSessionRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Session")).Return(nil)

	// Execute
	response, err := authUseCase.Login(context.Background(), login, "test-agent", "127.0.0.1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "access_token", response.AccessToken)
	assert.Equal(t, "refresh_token", response.RefreshToken)
	assert.NotNil(t, response.User)
	assert.Equal(t, user.Email, response.User.Email)

	// Verify all expectations were met
	mockUserRepo.AssertExpectations(t)
	mockSessionRepo.AssertExpectations(t)
	mockJWTService.AssertExpectations(t)
	mockPasswordService.AssertExpectations(t)
}

func TestAuthUseCase_GetProfile(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)

	authUseCase := NewAuthUseCase(mockUserRepo, mockSessionRepo, mockJWTService, mockPasswordService)

	userID := uuid.New()
	user := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		Password:  "hashed_password",
		FirstName: "John",
		LastName:  "Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Expectations
	mockUserRepo.On("GetByID", mock.Anything, userID).Return(user, nil)

	// Execute
	response, err := authUseCase.GetProfile(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotNil(t, response.User)
	assert.Equal(t, userID, response.User.ID)
	assert.Equal(t, user.Email, response.User.Email)
	assert.Equal(t, user.FirstName, response.User.FirstName)
	assert.Equal(t, user.LastName, response.User.LastName)

	// Verify all expectations were met
	mockUserRepo.AssertExpectations(t)
}

func TestAuthUseCase_GetProfile_UserNotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockSessionRepo := new(MockSessionRepository)
	mockJWTService := new(MockJWTService)
	mockPasswordService := new(MockPasswordService)

	authUseCase := NewAuthUseCase(mockUserRepo, mockSessionRepo, mockJWTService, mockPasswordService)

	userID := uuid.New()

	// Expectations
	mockUserRepo.On("GetByID", mock.Anything, userID).Return((*domain.User)(nil), errors.New("user not found"))

	// Execute
	response, err := authUseCase.GetProfile(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "user not found", err.Error())

	// Verify all expectations were met
	mockUserRepo.AssertExpectations(t)
}
