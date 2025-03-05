package main

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help information",
	Run: func(cmd *cobra.Command, args []string) {
		showHelp()
	},
}

func showHelp() {
	color.Yellow(`Available commands:

	help                           - show the help commands
	down                           - put the server into maintenance mode
	up                             - take the server out of maintenance mode
	version                        - print application version
	migrate                        - runs all up migrations that have not been run previously
	migrate down                   - reverses the most recent migration
	migrate reset                  - runs all down migrations in reverse order, and then all up migrations
	make migration <name> <format> - creates two new up and down migrations in the migrations folder; format=sql/fizz (default fizz)
	make auth                      - creates and runs migrations for authentication tables, and creates models and middleware
	make handler <name>            - creates a stub handler in the handlers directory
	make model <name>              - creates a new model in the data directory
	make session                   - creates a table in the database as a session store
	make mail <name>               - creates two starter mail templates in the mail directory
	
	`)
}

