package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scheduly.io/core/cmd/rest/handlers"
	"scheduly.io/core/cmd/rest/middleware"
	"scheduly.io/core/cmd/rest/routes"
	pg_repos "scheduly.io/core/internal/adapters/postgres/repositories"
	"scheduly.io/core/internal/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"scheduly.io/core/internal/config"
	"scheduly.io/core/internal/services"
)

type RestApp struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewRestApp(cfg *config.Config, db *gorm.DB) *RestApp {
	return &RestApp{
		cfg: cfg,
		db:  db,
	}
}

func (app *RestApp) Run() error {
	// Initialize repositories
	userRepo := pg_repos.NewUserRepository(app.db)
	sessionRepo := pg_repos.NewSessionRepository(app.db)

	// Initialize services
	jwtService := services.NewJWTService(
		app.cfg.JWT.AccessTokenSecret,
		app.cfg.JWT.RefreshTokenSecret,
		app.cfg.JWT.AccessTokenExpiry,
		app.cfg.JWT.RefreshTokenExpiry,
	)
	passwordService := services.NewPasswordService()

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(userRepo, sessionRepo, jwtService, passwordService, app.cfg.JWT)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupAuthRoutes(router, authHandler, authMiddleware)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Create the server
	server := &http.Server{
		Addr:         ":" + app.cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  app.cfg.Server.ReadTimeout,
		WriteTimeout: app.cfg.Server.WriteTimeout,
		IdleTimeout:  app.cfg.Server.IdleTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Server starting on port %s", app.cfg.Server.Port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("Server failed to start: %v", err)

	case sig := <-shutdown:
		log.Printf("Server is shutting down due to %v signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in 10s: %v", err)
			if err := server.Close(); err != nil {
				log.Fatalf("Could not stop server: %v", err)
			}
		}
	}

	return nil
}
