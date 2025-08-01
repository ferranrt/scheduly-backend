package helpers

import (
	"fmt"

	"ferranrt.com/scheduly-backend/internal/adapters/postgres"
	"ferranrt.com/scheduly-backend/internal/config"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (*gorm.DB, error) {
	cfg, err := config.New()
	fmt.Println(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	db, err := postgres.NewGormPostgreSQL(postgres.GormPostgreSQLConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
