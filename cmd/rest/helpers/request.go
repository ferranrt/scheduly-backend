package helpers

import (
	"strings"

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
