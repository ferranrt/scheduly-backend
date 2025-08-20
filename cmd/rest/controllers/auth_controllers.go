package controllers

import (
	"net/http"
	"strings"

	"bifur.app/core/cmd/rest/constants"
	"bifur.app/core/cmd/rest/helpers"
	"bifur.app/core/internal/domain"
	"bifur.app/core/internal/exceptions"
	"bifur.app/core/internal/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterController(ctx *gin.Context, authService ports.AuthService) {
	var registration domain.UserRegisterInput
	if err := ctx.ShouldBindJSON(&registration); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse(domain.ErrBadRequestPayload.Error()))
		return
	}

	requestMetadata := helpers.GetRequestMetadata(ctx)

	response, err := authService.Register(ctx.Request.Context(), &registration, requestMetadata.UserAgent, requestMetadata.IPAddress)
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

func RefreshTokenController(ctx *gin.Context, authService ports.AuthService) {
	var request domain.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse("Invalid request body"))
		return
	}

	response, err := authService.RefreshToken(ctx.Request.Context(), request.RefreshToken)
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

func GetProfileController(ctx *gin.Context, authService ports.AuthService) {

	userID := ctx.GetString(constants.UserIDClaimKey)

	userIDAsUUID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("User not authenticated"))
		return
	}

	profile, err := authService.GetProfile(ctx.Request.Context(), userIDAsUUID)
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

func LoginController(ctx *gin.Context, authService ports.AuthService) {
	var login domain.UserLoginInput
	if err := ctx.ShouldBindJSON(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse("Invalid request body"))
		return
	}

	meta := helpers.GetRequestMetadata(ctx)

	response, err := authService.Login(ctx.Request.Context(), &login, meta.UserAgent, meta.IPAddress)
	if err != nil {
		if err.Error() == exceptions.ErrAuthInvalidCredentials.Error() {
			ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse(err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to login"))
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func LogoutController(ctx *gin.Context, authService ports.AuthService) {
	var request domain.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse("Invalid request body"))
		return
	}

	err := authService.Logout(ctx.Request.Context(), request.RefreshToken)
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

// LogoutAllController handles user logout from all devices
func LogoutAllController(ctx *gin.Context, authService ports.AuthService) {
	userCtx, err := helpers.GetUserIdFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to logout from all devices"))
		return
	}

	if userCtx.UserID == "" {
		ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("User not authenticated"))
		return
	}

	err = authService.LogoutAll(ctx.Request.Context(), userCtx.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to logout from all devices"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out from all devices"})
}
