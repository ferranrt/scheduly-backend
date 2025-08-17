package helpers

import (
	"fmt"

	"buke.io/core/internal/adapters/postgres"
	"buke.io/core/internal/config"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (*gorm.DB, error) {
	cfg := config.New()

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
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
