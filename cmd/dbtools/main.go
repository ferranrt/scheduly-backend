package main

import (
	"fmt"
	"os"

	"scheduly.io/core/cmd/dbtools/commands"
)

var rootCmd = commands.RootCommand()

func init() {
	rootCmd.AddCommand(commands.Migrate())
	rootCmd.AddCommand(commands.Cleanup())
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
