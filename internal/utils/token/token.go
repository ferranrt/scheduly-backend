package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"bifur.app/core/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Common auth errors
var (
	ErrAuthHeaderMissing = errors.New("authentication required")
	ErrInvalidAuthFormat = errors.New("authorization header format must be bearer {token}")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

// JWTClaims holds the standard JWT claims plus our custom claims

// ValidateToken validates a JWT token string and returns the claims
func ValidateToken(tokenString string, secret []byte) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID uuid.UUID, email string, secret []byte, expiration time.Duration, sourceID string) (string, error) {
	claims := domain.JWTClaims{
		UserID:   userID,
		Email:    email,
		SourceID: sourceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ExtractToken extracts a token from query parameters or authorization header
func ExtractToken(c *gin.Context) (string, error) {
	// First try to get token from query parameter (common for WebSocket connections)
	token := c.Query("token")

	// If not in query, try header (for REST API)
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			return "", ErrAuthHeaderMissing
		}

		// Extract token from Bearer schema
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return "", ErrInvalidAuthFormat
		}
		token = parts[1]
	}

	return token, nil
}

// ExtractAndValidateToken combines extraction and validation
func ExtractAndValidateToken(c *gin.Context, secret []byte) (*domain.JWTClaims, error) {
	tokenString, err := ExtractToken(c)
	if err != nil {
		return nil, err
	}

	return ValidateToken(tokenString, secret)
}
