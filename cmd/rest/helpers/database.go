package helpers

import (
	"log"

	"gorm.io/gorm"
	"scheduly.io/core/internal/adapters/postgres"
	"scheduly.io/core/internal/config"
)

func GetDatabaseFromConfig(cfg *config.Config) (*gorm.DB, error) {
	db, err := postgres.NewGormPostgreSQL(postgres.GormPostgreSQLConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db, nil
}
