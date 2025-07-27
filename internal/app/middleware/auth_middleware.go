package middleware

import (
	"net/http"
	"strings"

	"ferranrt.com/scheduly-backend/internal/app/helpers"
	"ferranrt.com/scheduly-backend/internal/domain"
	"ferranrt.com/scheduly-backend/internal/ports/usecases"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authUseCase usecases.AuthUseCase
}

func NewAuthMiddleware(authUseCase usecases.AuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		authUseCase: authUseCase,
	}
}

func injectClaimsToContext(ctx *gin.Context, claims *domain.JWTClaims) {
	ctx.Set("user_id", claims.UserID)
	ctx.Set("email", claims.Email)
	ctx.Set("user_claims", claims)
}

// Authenticate middleware validates JWT tokens and sets user context
func (middleware *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("Authorization header required"))
			ctx.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("Invalid authorization header format"))
			ctx.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := middleware.authUseCase.ValidateToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("Invalid or expired token"))
			ctx.Abort()
			return
		}

		// Set user information in context
		injectClaimsToContext(ctx, claims)

		ctx.Next()
	}
}

// OptionalAuth middleware validates JWT tokens if present but doesn't require them
func (middleware *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Next()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the token
		claims, err := middleware.authUseCase.ValidateToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.Next()
			return
		}

		// Set user information in context
		injectClaimsToContext(ctx, claims)

		ctx.Next()
	}
}

// GetUserClaims extracts user claims from the context
func GetUserClaims(ctx *gin.Context) (*domain.JWTClaims, bool) {
	claims, exists := ctx.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*domain.JWTClaims)
	return userClaims, ok
}

// GetUserID extracts user ID from the context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetUserEmail extracts user email from the context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}

	userEmail, ok := email.(string)
	return userEmail, ok
}
