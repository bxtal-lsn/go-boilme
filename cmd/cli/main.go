package main

import (
	"errors"
	"os"

	"github.com/bxtal-lsn/go-boilme"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
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

// initCommands adds all child commands to the root command
func initCommands() {
	// Add all commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(helpCmd)
	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(makeCmd)
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

func validateInput() (string, string, string, string, error) {
	var arg1, arg2, arg3, arg4 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}

		if len(os.Args) >= 5 {
			arg4 = os.Args[4]
		}
	} else {
		color.Red("Error: command required")
		showHelp()
		return "", "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, arg4, nil
}

func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			exitGracefully(err)
		}

		path, err := os.Getwd()
		if err != nil {
			exitGracefully(err)
		}

		boil.RootPath = path
		boil.DB.DataType = os.Getenv("DATABASE_TYPE")
	}
}

