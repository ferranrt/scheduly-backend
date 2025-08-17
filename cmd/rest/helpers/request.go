package helpers

import (
	"log"
	"net/http"
	"strings"

	"buke.io/core/cmd/rest/constants"
	"buke.io/core/internal/domain"
	"github.com/gin-gonic/gin"
)

func GetClientIPFromRequest(ctx *gin.Context) string {

	if forwardedFor := ctx.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check for X-Real-IP header
	if realIP := ctx.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fallback to remote address
	return ctx.ClientIP()
}

func GetUserAgentFromRequest(ctx *gin.Context) string {
	return ctx.GetHeader("User-Agent")
}

func GetUserIdFromRequest(ctx *gin.Context) string {
	userID := ctx.GetString(constants.UserIDClaimKey)
	log.Println("GetUserIdFromRequest", userID)
	return userID
}

type RequestMetadata struct {
	UserAgent string
	IPAddress string
}

func GetRequestMetadata(ctx *gin.Context) RequestMetadata {
	return RequestMetadata{
		UserAgent: ctx.GetHeader("User-Agent"),
		IPAddress: GetClientIPFromRequest(ctx),
	}
}

func AbortUnauthorizedRequest(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusUnauthorized, BuildErrorResponse(err.Error()))
	ctx.Abort()
}

func AttachClaimsToContext(ctx *gin.Context, claims *domain.JWTClaims) {
	ctx.Set(constants.UserIDClaimKey, claims.UserID.String())
	ctx.Set(constants.EmailClaimKey, claims.Email)
	ctx.Set(constants.SourceIDClaimKey, claims.SourceID)
}
