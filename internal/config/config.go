package config

import (
	"log"
	"os"
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

func getEnv2(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("%s not set, defaulting to %s", key, defaultValue)
	return defaultValue
}

/* func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Invalid integer value for %s, defaulting to %d", key, defaultValue)
	}
	return defaultValue
} */

func getJWTDuration(specificKey string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(specificKey); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Invalid duration value for %s, defaulting to %d", specificKey, defaultValue)
	}

	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {

	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Invalid duration value for %s, defaulting to %d", key, defaultValue)
	}
	return defaultValue
}

func New() *Config {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv2("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv2("DB_HOST", "localhost"),
			Port:     getEnv2("DB_PORT", "5432"),
			User:     getEnv2("DB_USER", "postgres"),
			Password: getEnv2("DB_PASSWORD", "postgres"),
			DBName:   getEnv2("DB_NAME", "scheduly"),
			SSLMode:  getEnv2("DB_SSL_MODE", "disable"),
		},
		JWT: domain.JWTConfig{
			AtkSecret: getEnv2("JWT_ATK_SECRET_KEY", "your_access_token_secret_key_here"),
			RtkSecret: getEnv2("JWT_RTK_SECRET_KEY", "your_refresh_token_secret_key_here"),
			Expiry:    getJWTDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute),
		},
	}
	return config
}

func Print(cfg Config) {
	log.Printf("App Port: %s\n", cfg.Server.Port)
	log.Printf("DB Host: %s\n", cfg.Database.Host)
	log.Printf("DB Port: %s\n", cfg.Database.Port)
	log.Printf("DB Name: %s\n", cfg.Database.DBName)
	log.Printf("DB User: %s\n", cfg.Database.User)
	log.Printf("DB Password: %s\n", cfg.Database.Password)
	log.Printf("JWT Access Token Secret: %s\n", cfg.JWT.AtkSecret)
	log.Printf("JWT Refresh Token Secret: %s\n", cfg.JWT.RtkSecret)
	log.Printf("JWT Access Token Expiry: %s\n", cfg.JWT.Expiry)
}
