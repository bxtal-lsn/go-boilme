package main

import (
	"os"

	"github.com/bxtal-lsn/go-boilme"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

var boil boilme.Boilme

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "boilme",
	Short: "A powerful web application framework for Go",
	Long: `Boilme is a powerful web application framework for Go.
It provides tools for database migrations, authentication, and more.`,
}

func main() {
	// Initialize cobra commands
	initCommands()

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		exitGracefully(err)
	}
}

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Finished!")
	}

	os.Exit(0)
}

// initCommands adds all child commands to the root command
func initCommands() {
	// Add version command
	rootCmd.AddCommand(versionCmd)
}
