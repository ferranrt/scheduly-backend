package commands

import "github.com/spf13/cobra"

func RootCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "dbtools",
		Version: "1.0.0",
		Short:   "Database management tools for Scheduly backend",
		Long:    `A CLI tool for managing the Scheduly backend database operations like migrations and cleanup.`,
	}
}
