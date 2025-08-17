package config

import (
	"log"
	"os"
	"time"

	"buke.io/core/internal/domain"
	_ "github.com/joho/godotenv/autoload"
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
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	SSLMode    string
	LogEnabled bool
}

func getEnvVariable(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("%s not set, defaulting to %s", key, defaultValue)
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true"
	}
	log.Printf("%s not set, defaulting to %t", key, defaultValue)
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
			Port:         getEnvVariable("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:       getEnvVariable("DB_HOST", "localhost"),
			Port:       getEnvVariable("DB_PORT", "5432"),
			User:       getEnvVariable("DB_USER", "postgres"),
			Password:   getEnvVariable("DB_PASSWORD", "postgres"),
			DBName:     getEnvVariable("DB_NAME", "scheduly"),
			SSLMode:    getEnvVariable("DB_SSL_MODE", "disable"),
			LogEnabled: getEnvAsBool("DB_LOG_ENABLED", false),
		},
		JWT: domain.JWTConfig{
			AtkSecret: getEnvVariable("JWT_ATK_SECRET_KEY", "your_access_token_secret_key_here"),
			RtkSecret: getEnvVariable("JWT_RTK_SECRET_KEY", "your_refresh_token_secret_key_here"),
			Expiry:    getJWTDuration("JWT_ACCESS_TOKEN_DURATION", 15*time.Minute),
		},
	}
	return config
}

func Print(cfg Config) {
	log.Printf("--------------------------------")
	log.Printf("-------APP CONFIG---------------")
	log.Printf("--------------------------------")
	log.Printf("App Port: %s\n", cfg.Server.Port)
	log.Printf("--------------------------------")
	log.Printf("-------DB CONFIG---------------")
	log.Printf("--------------------------------")
	log.Printf("DB Host: %s\n", cfg.Database.Host)
	log.Printf("DB Port: %s\n", cfg.Database.Port)
	log.Printf("DB Name: %s\n", cfg.Database.DBName)
	log.Printf("DB User: %s\n", cfg.Database.User)
	log.Printf("DB Password: %s\n", cfg.Database.Password)
	log.Printf("DB Log Enabled: %t\n", cfg.Database.LogEnabled)
	log.Printf("--------------------------------")
	log.Printf("-------JWT CONFIG---------------")
	log.Printf("--------------------------------")
	log.Printf("JWT ATK Secret: %s\n", cfg.JWT.AtkSecret)
	log.Printf("JWT RTK Secret: %s\n", cfg.JWT.RtkSecret)
	log.Printf("JWT Expiry: %s\n", cfg.JWT.Expiry)
	log.Printf("--------------------------------")
}
