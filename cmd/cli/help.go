package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// helpCmd represents an enhanced help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show help information",
	Run: func(cmd *cobra.Command, args []string) {
		showHelp()
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}

func showHelp() {
	// Create a beautiful header
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println("Boilme CLI")

	pterm.Info.Println("A powerful web application framework for Go")
	pterm.Println()

	// Display available commands
	tableData := pterm.TableData{
		{"Command", "Description"},
		{"help", "Show help information"},
		{"version", "Print application version"},
		{"new [name]", "Create a new application"},
		{"up", "Take the server out of maintenance mode"},
		{"down", "Put the server into maintenance mode"},
		{"migrate", "Run database migrations"},
		{"migrate down", "Reverse the most recent migration"},
		{"migrate reset", "Run all down migrations, then all up migrations"},
	}

	// Create section for make commands
	pterm.DefaultSection.Println("Make Commands")

	makeTableData := pterm.TableData{
		{"Command", "Description"},
		{"make migration [name] [format]", "Create migration files (format=sql/fizz)"},
		{"make auth", "Generate authentication system"},
		{"make handler [name]", "Create a handler in the handlers directory"},
		{"make model [name]", "Create a model in the data directory"},
		{"make session", "Create a database session store"},
		{"make mail [name]", "Create mail templates"},
		{"make key", "Generate a random encryption key"},
	}

	// Print the tables
	pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	pterm.Println()
	pterm.DefaultTable.WithHasHeader().WithData(makeTableData).Render()

	// Add usage examples
	pterm.DefaultSection.Println("Examples")

	// Simple alternative using styled text
	cyan := pterm.FgCyan.Sprint("•")
	gray := pterm.FgGray.Sprint("•")

	pterm.Printf("%s %s\n", cyan, pterm.FgLightCyan.Sprint("Create a new application: boilme new myapp"))
	pterm.Printf("  %s %s\n", gray, pterm.FgGray.Sprint("Create a migration: boilme make migration create_users_table"))
	pterm.Printf("%s %s\n", cyan, pterm.FgLightCyan.Sprint("Run migrations: boilme migrate"))
}

