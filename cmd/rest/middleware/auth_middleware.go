package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"scheduly.io/core/cmd/rest/constants"
	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/exceptions"
	"scheduly.io/core/internal/ports"
	"scheduly.io/core/internal/utils/token"
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

func injectClaimsToContext(ctx *gin.Context, claims *domain.JWTClaims) {
	ctx.Set(constants.UserIDClaimKey, claims.UserID.String())
	ctx.Set(constants.EmailClaimKey, claims.Email)
	ctx.Set(constants.SourceIDClaimKey, claims.SourceID)
}

func abortWithUnauthorized(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse(err.Error()))
	ctx.Abort()
}

// Authenticate middleware validates JWT tokens and sets user context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tk, err := token.ExtractToken(ctx)
		if err != nil {
			abortWithUnauthorized(ctx, err)
			return
		}
		decoded, err := token.ExtractAndValidateToken(ctx, []byte(m.jwtConfig.AtkSecret))

		if err != nil {
			abortWithUnauthorized(ctx, err)
			return
		}

		source, err := m.sourceRepo.GetByID(ctx.Request.Context(), decoded.SourceID)
		if err != nil {
			abortWithUnauthorized(ctx, err)
			return
		}
		if !source.IsActive || time.Now().After(source.RefreshTokenExpiresAt) {
			abortWithUnauthorized(ctx, exceptions.ErrAuthSourceExpiredOrInvalid)
			return
		}

		// Validate the token
		claims, err := m.authService.ValidateToken(ctx.Request.Context(), tk)
		if err != nil {
			abortWithUnauthorized(ctx, err)
			return
		}

		// Set user information in context

		injectClaimsToContext(ctx, claims)

		ctx.Next()
	}
}
