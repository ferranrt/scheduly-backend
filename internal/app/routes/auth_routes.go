package routes

import (
	"ferranrt.com/scheduly-backend/internal/app/handlers"
	"ferranrt.com/scheduly-backend/internal/app/middleware"
	"github.com/gin-gonic/gin"
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
