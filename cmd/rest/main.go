package main

import (
	"log"

	"ferranrt.com/scheduly-backend/internal/adapters/postgres"
	"ferranrt.com/scheduly-backend/internal/app"
	"ferranrt.com/scheduly-backend/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
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
		log.Fatalf("Failed to connect to database: %v", err)
	}

	appInstance := app.NewRestApp(cfg, db)

	if err := appInstance.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
