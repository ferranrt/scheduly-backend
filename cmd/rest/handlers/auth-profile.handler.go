package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"scheduly.io/core/cmd/rest/helpers"
)

// GetProfile handles user profile retrieval
func (h *AuthRootHandler) GetProfile(ctx *gin.Context) {

	userID := helpers.GetUserIdFromRequest(ctx)
	fmt.Println(userID)
	userIDAsUUID, err := uuid.Parse(userID)
	fmt.Println(userIDAsUUID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("User not authenticated"))
		return
	}

	profile, err := h.authUseCase.GetProfile(ctx.Request.Context(), userIDAsUUID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, helpers.BuildErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to retrieve user profile"))
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
