package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"scheduly.io/core/cmd/dbtools/helpers"
	"scheduly.io/core/internal/adapters/postgres/dbmodels"
	"scheduly.io/core/internal/adapters/postgres/migrations"
)

func Rebuild() *cobra.Command {
	return &cobra.Command{
		Use:   "rebuild",
		Short: "Rebuild database tables",
		Long:  `Drop all tables and recreate them from scratch. Use with caution!`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runRebuild(); err != nil {
				log.Fatalf("Cleanup failed: %v", err)
			}
			fmt.Println("Database cleanup completed successfully!")
		},
	}
}

func runRebuild() error {
	db, err := helpers.GetDatabaseConnection()
	if err != nil {
		return err
	}

	// Drop all tables and recreate them
	models := []interface{}{
		&dbmodels.User{},
		&dbmodels.Source{},
	}

	// Drop all tables
	if err := db.Migrator().DropTable(models...); err != nil {
		return fmt.Errorf("failed to drop tables: %v", err)
	}

	// Recreate tables
	if err := migrations.Migrate(db); err != nil {
		return fmt.Errorf("failed to recreate tables: %v", err)
	}

	return nil
}
