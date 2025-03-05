package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const version = "1.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the current version of Boilme.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow("Boilme Framework version: %s", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
