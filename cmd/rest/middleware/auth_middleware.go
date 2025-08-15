package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"scheduly.io/core/cmd/rest/constants"
	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/ports"
	"scheduly.io/core/internal/utils/token"
)

type AuthMiddleware struct {
	authService ports.AuthService
	jwtConfig   domain.JWTConfig
}

func NewAuthMiddleware(authUseCase ports.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authUseCase,
	}
}

func injectClaimsToContext(ctx *gin.Context, claims *domain.JWTClaims) {
	log.Println("injectClaimsToContext", claims)
	ctx.Set(constants.UserIDClaimKey, claims.UserID.String())
	ctx.Set(constants.EmailClaimKey, claims.Email)
}

// Authenticate middleware validates JWT tokens and sets user context
func (middleware *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tk, err := token.ExtractToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse(err.Error()))
			ctx.Abort()
			return
		}

		test, err := token.ExtractAndValidateToken(ctx, []byte(middleware.jwtConfig.AtkSecret))
		fmt.Println("test", test)
		fmt.Println("err", err)

		// Validate the token
		claims, err := middleware.authService.ValidateToken(ctx.Request.Context(), tk)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("Invalid or expired token"))
			ctx.Abort()
			return
		}

		// Set user information in context

		fmt.Println("claims", claims)
		injectClaimsToContext(ctx, claims)

		ctx.Next()
	}
}
