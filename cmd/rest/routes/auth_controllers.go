package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"scheduly.io/core/cmd/rest/constants"
	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/internal/domain"
	"scheduly.io/core/internal/ports"
)

func RegisterController(ctx *gin.Context, authService ports.AuthService) {
	var registration domain.UserRegistrationInput
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

// GetProfileController handles user profile retrieval
func GetProfileController(ctx *gin.Context, authService ports.AuthService) {
	helpers.PrintContextInternals(ctx, false)
	userID := ctx.GetString(constants.UserIDClaimKey)
	log.Println("GetProfile -> userID", userID)
	userIDAsUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Println("GetProfileError", err)
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

	requestMetadata := helpers.GetRequestMetadata(ctx)

	response, err := authService.Login(ctx.Request.Context(), &login, requestMetadata.UserAgent, requestMetadata.IPAddress)
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
	userID := helpers.GetUserIdFromRequest(ctx)

	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, helpers.BuildErrorResponse("User not authenticated"))
		return
	}

	err := authService.LogoutAll(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.BuildErrorResponse("Failed to logout from all devices"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out from all devices"})
}
