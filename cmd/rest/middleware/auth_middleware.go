package middleware

import (
	"time"

	"bifur.app/core/cmd/rest/helpers"
	"bifur.app/core/internal/domain"
	"bifur.app/core/internal/exceptions"
	"bifur.app/core/internal/ports"
	"bifur.app/core/internal/utils/token"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService ports.AuthService
	sourceRepo  ports.SourceRepository
	jwtConfig   domain.JWTConfig
}

func NewAuthMiddleware(authUseCase ports.AuthService, sourceRepo ports.SourceRepository, jwtConfig domain.JWTConfig) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authUseCase,
		sourceRepo:  sourceRepo,
		jwtConfig:   jwtConfig,
	}
}

// Authenticate middleware validates JWT tokens and sets user context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tk, err := token.ExtractToken(ctx)
		if err != nil {
			helpers.AbortUnauthorizedRequest(ctx, err)
			return
		}
		decoded, err := token.ExtractAndValidateToken(ctx, []byte(m.jwtConfig.AtkSecret))

		if err != nil {
			helpers.AbortUnauthorizedRequest(ctx, err)
			return
		}

		source, err := m.sourceRepo.GetByID(ctx.Request.Context(), decoded.SourceID)
		if err != nil {
			helpers.AbortUnauthorizedRequest(ctx, err)
			return
		}
		if !source.IsActive || time.Now().After(source.RefreshTokenExpiresAt) {
			helpers.AbortUnauthorizedRequest(ctx, exceptions.ErrAuthSourceExpiredOrInvalid)
			return
		}

		// Validate the token
		claims, err := m.authService.ValidateToken(ctx.Request.Context(), tk)
		if err != nil {
			helpers.AbortUnauthorizedRequest(ctx, err)
			return
		}

		// Set user information in context

		helpers.AttachClaimsToContext(ctx, claims)

		ctx.Next()
	}
}
