package main

import (
	"log"

	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/cmd/rest/server"
	"scheduly.io/core/internal/config"
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
