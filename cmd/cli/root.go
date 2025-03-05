package cmd

import (
	"fmt"
	"os"

	"github.com/bxtal-lsn/go-boilme"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var boil boilme.Boilme

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "boilme",
	Short: "A powerful web application framework for Go",
	Long: `Boilme is a powerful web application framework for Go, 
designed to make web development faster and more enjoyable.

It provides tools for database migrations, authentication, 
file handling, caching, and more.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		exitGracefully(err)
	}
}

func init() {
	rootPath, _ := os.Getwd()
	err := boil.New(rootPath)
	if err != nil {
		fmt.Printf("Error initializing boilme: %v\n", err)
		os.Exit(1)
	}
	// Here you will define your flags and configuration settings
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is .env)")

	// Add completion command
	rootCmd.AddCommand(completionCmd)
}

// exitGracefully exits the application with an optional error message
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

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:
$ source <(boilme completion bash)

Zsh:
$ source <(boilme completion zsh)

fish:
$ boilme completion fish | source

PowerShell:
PS> boilme completion powershell | Out-String | Invoke-Expression
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

// setup initializes the application
func setupRoot(cmd *cobra.Command, args []string) { // Load .env file
	// This should be done in each command that needs it
	if cmd.Use != "new" && cmd.Use != "version" && cmd.Use != "help" {
		err := os.Chdir(boil.RootPath)
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
