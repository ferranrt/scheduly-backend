package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/internal/domain"
)

// RefreshToken handles token refresh
func (h *AuthRootHandler) RefreshToken(ctx *gin.Context) {
	var request domain.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse("Invalid request body"))
		return
	}

	response, err := h.authUseCase.RefreshToken(ctx.Request.Context(), request.RefreshToken)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "expired") {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to refresh token"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
