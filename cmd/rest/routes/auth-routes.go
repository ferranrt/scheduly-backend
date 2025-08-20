package routes

import (
	"bifur.app/core/cmd/rest/controllers"
	"bifur.app/core/internal/ports"
	"github.com/gin-gonic/gin"
)

type AuthRoutesDeps struct {
	AuthService ports.AuthService
}

func SetupPublicAuthRoutes(router *gin.RouterGroup, deps *AuthRoutesDeps) {
	router.POST("/register", func(ctx *gin.Context) { controllers.RegisterController(ctx, deps.AuthService) })
	router.POST("/login", func(ctx *gin.Context) { controllers.LoginController(ctx, deps.AuthService) })
	router.POST("/refresh", func(ctx *gin.Context) { controllers.RefreshTokenController(ctx, deps.AuthService) })
	router.POST("/logout", func(ctx *gin.Context) { controllers.LogoutController(ctx, deps.AuthService) })
}

func SetupProtectedAuthRoutes(router *gin.RouterGroup, deps *AuthRoutesDeps) {

	router.GET("/profile", func(ctx *gin.Context) { controllers.GetProfileController(ctx, deps.AuthService) })
	router.POST("/logout-all", func(ctx *gin.Context) { controllers.LogoutAllController(ctx, deps.AuthService) })
}
