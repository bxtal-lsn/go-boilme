package main

import (
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down|reset]",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		direction := "up"
		count := ""

		if len(args) > 0 {
			direction = args[0]
		}

		if len(args) > 1 {
			count = args[1]
		}

		err := doMigrate(direction, count)
		if err != nil {
			exitGracefully(err)
		}

		exitGracefully(nil, "Migrations complete!")
	},
}

func doMigrate(arg2, arg3 string) error {
	// dsn := getDSN()
	checkForDB()

	tx, err := boil.PopConnect()
	if err != nil {
		exitGracefully(err)
	}
	defer tx.Close()

	// run the migration command
	switch arg2 {
	case "up":
		err := boil.RunPopMigrations(tx)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
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
		showHelp()
	}

	return nil
}
