package commands

import (
	"fmt"
	"log"

	"buke.io/core/cmd/dbtools/helpers"
	"buke.io/core/internal/adapters/postgres/migrations"
	"github.com/spf13/cobra"
)

func Migrate() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  `Execute database migrations to create or update database schema.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runMigrations(); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Println("Database migrations completed successfully!")
		},
	}
}

func runMigrations() error {
	db, err := helpers.GetDatabaseConnection()
	if err != nil {
		return err
	}

	return migrations.Migrate(db)
}
