package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/internal/domain"
)

// Register handles user registration
func (h *AuthRootHandler) Register(ctx *gin.Context) {
	var registration domain.UserRegistrationInput
	if err := ctx.ShouldBindJSON(&registration); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse(domain.ErrBadRequestPayload.Error()))
		return
	}

	requestMetadata := helpers.GetRequestMetadata(ctx)

	response, err := h.authUseCase.Register(ctx.Request.Context(), &registration, requestMetadata.UserAgent, requestMetadata.IPAddress)
	if err != nil {
		if err.Error() == "user already exists" {
			ctx.JSON(http.StatusConflict, helpers.BuildErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to register user"))
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
