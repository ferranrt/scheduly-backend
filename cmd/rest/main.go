package main

import (
	"log"

	"bifur.app/core/cmd/rest/helpers"
	"bifur.app/core/cmd/rest/server"
	"bifur.app/core/internal/config"
)

func main() {
	log.Println("Starting server")

	log.Println("Getting configuration")
	cfg := config.New()
	config.Print(*cfg)

	log.Println("Connecting to database")
	db, err := helpers.GetDatabaseFromConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	appInstance := server.NewRestApp(cfg, db)

	if err := appInstance.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
