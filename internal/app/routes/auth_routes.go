package routes

import (
	"github.com/gin-gonic/gin"
	"scheduly.io/core/internal/app/handlers"
	"scheduly.io/core/internal/app/middleware"
)

func SetupAuthRoutes(router *gin.Engine, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	auth := router.Group("/api/v1/auth")
	{
		// Public routes
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)

		// Protected routes
		protected := auth.Group("")
		protected.Use(authMiddleware.Authenticate())
		{
			protected.GET("/profile", authHandler.GetProfile)
			protected.POST("/logout-all", authHandler.LogoutAll)
		}
	}
}
