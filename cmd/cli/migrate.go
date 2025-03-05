package cmd

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate [direction]",
	Short: "Run database migrations",
	Long: `Run database migrations in the specified direction:
  
  - up: Run all pending migrations (default)
  - down: Rollback one migration
  - reset: Roll back all migrations and run them again`,
	Run: func(cmd *cobra.Command, args []string) {
		direction := "up"
		if len(args) > 0 {
			direction = args[0]
		}

		all, _ := cmd.Flags().GetBool("all")

		// Start spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Running migrations..."
		s.Color("blue")
		s.Start()

		err := doMigrate(direction, all)

		s.Stop()
		if err != nil {
			exitGracefully(err)
		}

		color.Green("✓ Migrations completed successfully!")
	},
}

// migrateFreshCmd represents the migrate fresh command
var migrateFreshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Drop all tables and run migrations",
	Long:  `Drop all tables from the database and then run all migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Dropping tables and running migrations..."
		s.Color("blue")
		s.Start()

		// First migrate down all
		err := doMigrate("down", true)
		if err != nil {
			s.Stop()
			exitGracefully(err)
		}

		// Then migrate up
		err = doMigrate("up", false)

		s.Stop()
		if err != nil {
			exitGracefully(err)
		}

		color.Green("✓ Fresh migration completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateFreshCmd)

	// Add flags
	migrateCmd.Flags().BoolP("all", "a", false, "Apply to all migrations (for down command)")
}

func doMigrate(direction string, all bool) error {
	// Implementation adapted from your existing migrate.go
	checkForDB()

	tx, err := boil.PopConnect()
	if err != nil {
		return err
	}
	defer tx.Close()

	// run the migration command
	switch direction {
	case "up":
		err := boil.RunPopMigrations(tx)
		if err != nil {
			return err
		}

	case "down":
		if all {
			err := boil.PopMigrateDown(tx, -1)
			if err != nil {
				return err
			}
		} else {
			err := boil.PopMigrateDown(tx, 1)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := boil.PopMigrateReset(tx)
		if err != nil {
			return err
		}
	default:
		exitGracefully(nil, "Invalid migration direction")
	}

	return nil
}
