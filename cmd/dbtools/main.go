package main

import (
	"fmt"
	"os"

	"buke.io/core/cmd/dbtools/commands"
)

var rootCmd = commands.RootCommand()

func init() {
	rootCmd.AddCommand(commands.Migrate())
	rootCmd.AddCommand(commands.Rebuild())
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
