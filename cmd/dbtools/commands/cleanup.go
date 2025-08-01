package commands

import (
	"fmt"
	"log"

	"ferranrt.com/scheduly-backend/cmd/dbtools/helpers"
	"ferranrt.com/scheduly-backend/internal/adapters/postgres/dbmodels"
	"ferranrt.com/scheduly-backend/internal/adapters/postgres/migrations"
	"github.com/spf13/cobra"
)

func Cleanup() *cobra.Command {
	return &cobra.Command{
		Use:   "cleanup",
		Short: "Clean up database tables",
		Long:  `Drop all tables and recreate them from scratch. Use with caution!`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCleanup(); err != nil {
				log.Fatalf("Cleanup failed: %v", err)
			}
			fmt.Println("Database cleanup completed successfully!")
		},
	}
}

func runCleanup() error {
	db, err := helpers.GetDatabaseConnection()
	if err != nil {
		return err
	}

	// Drop all tables and recreate them
	models := []interface{}{
		&dbmodels.User{},
		&dbmodels.Session{},
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
