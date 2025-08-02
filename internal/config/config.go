package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"scheduly.io/core/internal/domain"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      domain.JWTConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var (
	once           sync.Once
	configInstance *Config
)

func New() *Config {
	once.Do(func() {
		config := &Config{
			Server: ServerConfig{
				Port:         getEnv("SERVER_PORT", "8080"),
				ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
				WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
				IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", "postgres"),
				DBName:   getEnv("DB_NAME", "scheduly"),
				SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			},
			JWT: domain.JWTConfig{
				AccessTokenSecret:  getEnv("JWT_SECRET_KEY", "your_access_token_secret_key_here"),
				RefreshTokenSecret: getEnv("JWT_SECRET_KEY", "your_refresh_token_secret_key_here"),
				AccessTokenExpiry:  getJWTDuration("JWT_ACCESS_TOKEN_DURATION", "JWT_DURATION", 15*time.Minute),
				RefreshTokenExpiry: getJWTDuration("JWT_REFRESH_TOKEN_DURATION", "JWT_DURATION", 30*24*time.Hour), // 30 days
			},
		}

		if err := config.validate(); err != nil {
			panic(fmt.Errorf("invalid configuration: %w", err))
		}

		configInstance = config
	})

	return configInstance
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}

	if c.JWT.AccessTokenSecret == "" {
		return fmt.Errorf("JWT access token secret key is required")
	}

	if c.JWT.RefreshTokenSecret == "" {
		return fmt.Errorf("JWT refresh token secret key is required")
	}

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getDurationEnv gets a duration from an environment variable or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {

	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getJWTDuration gets a JWT duration from environment variables with fallback support
// It first checks for a specific duration variable, then falls back to JWT_DURATION
func getJWTDuration(specificKey, fallbackKey string, defaultValue time.Duration) time.Duration {
	// First try the specific key (e.g., JWT_ACCESS_TOKEN_DURATION)
	if value, exists := os.LookupEnv(specificKey); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}

	// If specific key doesn't exist or is invalid, try the fallback key (JWT_DURATION)
	if value, exists := os.LookupEnv(fallbackKey); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}

	// Return default value if neither key exists or both are invalid
	return defaultValue
}
