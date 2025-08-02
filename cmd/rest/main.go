package main

import (
	"fmt"
	"log"

	"scheduly.io/core/cmd/rest/helpers"
	"scheduly.io/core/cmd/rest/server"
	"scheduly.io/core/internal/config"
)

func main() {
	cfg := config.New()
	fmt.Println(cfg)
	db, err := helpers.GetDatabaseFromConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	appInstance := server.NewRestApp(cfg, db)

	if err := appInstance.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
