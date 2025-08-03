package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/internal/domain"
)

// Register handles user registration

// Login handles user login
func (h *AuthRootHandler) Login(ctx *gin.Context) {
	var login domain.UserLoginInput
	if err := ctx.ShouldBindJSON(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse("Invalid request body"))
		return
	}

	requestMetadata := helpers.GetRequestMetadata(ctx)

	response, err := h.authUseCase.Login(ctx.Request.Context(), &login, requestMetadata.UserAgent, requestMetadata.IPAddress)
	if err != nil {
		if err.Error() == "invalid credentials" {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to login"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
