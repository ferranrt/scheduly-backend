package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/internal/domain"
)

// Logout handles user logout
func (h *AuthRootHandler) Logout(ctx *gin.Context) {
	var request domain.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse("Invalid request body"))
		return
	}

	err := h.authUseCase.Logout(ctx.Request.Context(), request.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to logout"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// LogoutAll handles user logout from all devices
func (h *AuthRootHandler) LogoutAll(ctx *gin.Context) {
	userID := helpers.GetUserIdFromRequest(ctx)

	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("User not authenticated"))
		return
	}

	err := h.authUseCase.LogoutAll(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to logout from all devices"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out from all devices"})
}
