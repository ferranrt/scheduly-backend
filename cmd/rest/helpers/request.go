package helpers

import (
	"net/http"
	"strings"

	"bifur.app/core/cmd/rest/constants"
	"bifur.app/core/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type UserCtx struct {
	UserID string
	AsUUID uuid.UUID
}

func GetUserIdFromRequest(ctx *gin.Context) (UserCtx, error) {
	userID := ctx.GetString(constants.UserIDClaimKey)
	userIDUUID, err := uuid.Parse(userID)
	if err != nil {
		return UserCtx{}, err
	}
	return UserCtx{UserID: userID, AsUUID: userIDUUID}, nil
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
