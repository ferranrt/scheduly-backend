package helpers

import (
	"log"

	"bifur.app/core/internal/adapters/postgres"
	"bifur.app/core/internal/config"
	"gorm.io/gorm"
)

func GetDatabaseFromConfig(cfg *config.Config) (*gorm.DB, error) {
	db, err := postgres.NewGormPostgreSQL(postgres.GormPostgreSQLConfig{
		Host:       cfg.Database.Host,
		Port:       cfg.Database.Port,
		User:       cfg.Database.User,
		Password:   cfg.Database.Password,
		DBName:     cfg.Database.DBName,
		SSLMode:    cfg.Database.SSLMode,
		LogEnabled: cfg.Database.LogEnabled,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db, nil
}
