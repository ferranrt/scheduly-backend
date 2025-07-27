package handlers

import (
	"net/http"
	"strings"

	"ferranrt.com/scheduly-backend/internal/app/helpers"
	"ferranrt.com/scheduly-backend/internal/domain"
	"ferranrt.com/scheduly-backend/internal/ports/usecases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authUseCase usecases.AuthUseCase
}

func NewAuthHandler(authUseCase usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(ctx *gin.Context) {
	var registration domain.UserRegistrationInput
	if err := ctx.ShouldBindJSON(&registration); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse(domain.ErrBadRequestPayload.Error()))
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	ipAddress := helpers.GetClientIPFromRequest(ctx)

	response, err := h.authUseCase.Register(ctx.Request.Context(), &registration, userAgent, ipAddress)
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

// Login handles user login
func (h *AuthHandler) Login(ctx *gin.Context) {
	var login domain.UserLoginInput
	if err := ctx.ShouldBindJSON(&login); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.BuildErrorResponse("Invalid request body"))
		return
	}

	userAgent := helpers.GetUserAgentFromRequest(ctx)
	ipAddress := helpers.GetClientIPFromRequest(ctx)

	response, err := h.authUseCase.Login(ctx.Request.Context(), &login, userAgent, ipAddress)
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
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
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

// Logout handles user logout
func (h *AuthHandler) Logout(ctx *gin.Context) {
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

// LogoutAll handles logout from all devices
func (h *AuthHandler) LogoutAll(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
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

// GetProfile handles user profile retrieval
func (h *AuthHandler) GetProfile(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	userIDAsUUID, err := uuid.Parse(userID)
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
