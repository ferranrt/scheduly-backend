package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"buke.io/core/cmd/rest/middleware"
	"buke.io/core/cmd/rest/routes"
	pg_repos "buke.io/core/internal/adapters/postgres/repositories"
	"buke.io/core/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"buke.io/core/internal/config"
)

type RestApp struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewRestApp(cfg *config.Config, db *gorm.DB) *RestApp {
	return &RestApp{cfg: cfg, db: db}
}

func createServer(cfg *config.Config, engine *gin.Engine) *http.Server {

	return &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func (app *RestApp) Run() error {
	// Initialize repositories
	userRepository := pg_repos.NewUserRepository(app.db)
	sourceRepository := pg_repos.NewSourceRepository(app.db)

	// Initialize services
	authService := services.NewAuthService(userRepository, sourceRepository, app.cfg.JWT)

	// Initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(authService, sourceRepository, app.cfg.JWT)

	// Setup Gin router
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware("*"))

	/* Public routes */
	publicGroup := router.Group("/api/v1")

	/* Protected routes */
	protectedGroup := router.Group("/api/v1")
	protectedGroup.Use(authMiddleware.Authenticate())

	// Setup routes
	// Health Routes
	routes.SetupHealthRoutes(router)
	// Auth Routes
	routes.SetupPublicAuthRoutes(publicGroup, &routes.AuthRoutesDeps{AuthService: authService})
	routes.SetupProtectedAuthRoutes(protectedGroup, &routes.AuthRoutesDeps{AuthService: authService})

	// Create the server
	server := createServer(app.cfg, router)

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
