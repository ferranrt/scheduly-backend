package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"scheduly.io/core/internal/domain"
)

type JWTService interface {
	GenerateAccessToken(user *domain.User) (string, error)
	GenerateRefreshToken(user *domain.User) (string, error)
	ValidateToken(tokenString string) (*domain.JWTClaims, error)
	GenerateRandomToken() (string, error)
}

type jwtService struct {
	config domain.JWTConfig
}

func NewJWTService(config domain.JWTConfig) JWTService {
	return &jwtService{
		config: config,
	}
}

func (s *jwtService) GenerateAccessToken(user *domain.User) (string, error) {
	claims := domain.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"type":    claims.Type,
		"exp":     time.Now().Add(s.config.AccessTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(s.config.AccessTokenSecret))
}

func (s *jwtService) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := domain.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Type:   "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"type":    claims.Type,
		"exp":     time.Now().Add(s.config.RefreshTokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(s.config.RefreshTokenSecret))
}

func (s *jwtService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	// First try to validate as access token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.AccessTokenSecret), nil
	})

	if err == nil && token.Valid {
		return s.extractClaims(token)
	}

	// If access token validation fails, try refresh token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.RefreshTokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return s.extractClaims(token)
}

func (s *jwtService) extractClaims(token *jwt.Token) (*domain.JWTClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}
	parsedUserID, err := uuid.Parse(userID)

	if err != nil {
		return nil, errors.New("invalid user_id format in token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in token")
	}

	tokenType, ok := claims["type"].(string)
	if !ok {
		return nil, errors.New("invalid token type")
	}

	return &domain.JWTClaims{
		UserID: parsedUserID,
		Email:  email,
		Type:   tokenType,
	}, nil
}

func (s *jwtService) GenerateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
