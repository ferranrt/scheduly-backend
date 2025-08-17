package routes

import (
	"buke.io/core/internal/ports"
	"github.com/gin-gonic/gin"
)

type AuthRoutesDeps struct {
	AuthService ports.AuthService
}

func SetupPublicAuthRoutes(router *gin.RouterGroup, deps *AuthRoutesDeps) {
	router.POST("/register", func(ctx *gin.Context) { RegisterController(ctx, deps.AuthService) })
	router.POST("/login", func(ctx *gin.Context) { LoginController(ctx, deps.AuthService) })
	router.POST("/refresh", func(ctx *gin.Context) { RefreshTokenController(ctx, deps.AuthService) })
	router.POST("/logout", func(ctx *gin.Context) { LogoutController(ctx, deps.AuthService) })
}

func SetupProtectedAuthRoutes(router *gin.RouterGroup, deps *AuthRoutesDeps) {

	router.GET("/profile", func(ctx *gin.Context) { GetProfileController(ctx, deps.AuthService) })
	router.POST("/logout-all", func(ctx *gin.Context) { LogoutAllController(ctx, deps.AuthService) })
}
