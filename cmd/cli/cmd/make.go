package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make [subcommand]",
	Short: "Generate various components",
	Long: `Generate various components for your Boilme application:
  
  - migration: Create a new database migration
  - model: Create a new model
  - handler: Create a new HTTP handler
  - auth: Generate auth system (tables, models, handlers)
  - mail: Create new mail templates
  - session: Create a new session table
  - key: Generate a random encryption key
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

// Subcommands
var makeMigrationCmd = &cobra.Command{
	Use:   "migration [name] [type]",
	Short: "Create a new migration",
	Long: `Create a new database migration file.
The type parameter is optional and can be 'fizz' (default) or 'sql'.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup necessary for database operations
		setup(cmd, args)

		name := args[0]
		migrationType := "fizz"
		if len(args) > 1 && args[1] == "sql" {
			migrationType = "sql"
		}

		// Generate migration
		err := doMakeMigration(name, migrationType)
		if err != nil {
			exitGracefully(err)
		}

		color.Green("Created migration %s", name)
	},
}

var makeModelCmd = &cobra.Command{
	Use:   "model [name]",
	Short: "Create a new model",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create model
		err := doMakeModel(args[0])
		if err != nil {
			exitGracefully(err)
		}

		color.Green("Created model %s", args[0])
	},
}

var makeHandlerCmd = &cobra.Command{
	Use:   "handler [name]",
	Short: "Create a new handler",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create handler
		err := doMakeHandler(args[0])
		if err != nil {
			exitGracefully(err)
		}

		color.Green("Created handler %s", args[0])
	},
}

var makeAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate auth system",
	Run: func(cmd *cobra.Command, args []string) {
		setup(cmd, args)
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}
	},
}

var makeMailCmd = &cobra.Command{
	Use:   "mail [name]",
	Short: "Create mail templates",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create mail templates
		err := doMakeMail(args[0])
		if err != nil {
			exitGracefully(err)
		}

		color.Green("Created mail templates for %s", args[0])
	},
}

var makeSessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Create session table",
	Run: func(cmd *cobra.Command, args []string) {
		setup(cmd, args)
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}

		color.Green("Created session table")
	},
}

var makeKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Generate a random encryption key",
	Run: func(cmd *cobra.Command, args []string) {
		rnd := boil.RandomString(32)
		color.Yellow("32 character encryption key: %s", rnd)
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)

	// Add subcommands to make
	makeCmd.AddCommand(makeMigrationCmd)
	makeCmd.AddCommand(makeModelCmd)
	makeCmd.AddCommand(makeHandlerCmd)
	makeCmd.AddCommand(makeAuthCmd)
	makeCmd.AddCommand(makeMailCmd)
	makeCmd.AddCommand(makeSessionCmd)
	makeCmd.AddCommand(makeKeyCmd)
}

// Functions that implement the actual functionality
// These would be similar to the functions in your existing make.go file

func doMakeMigration(name, migrationType string) error {
	// Implementation from your existing code
	checkForDB()

	if name == "" {
		return errors.New("you must give the migration a name")
	}

	var up, down string

	if migrationType == "fizz" || migrationType == "" {
		upBytes, _ := templateFS.ReadFile("templates/migrations/migration_up.fizz")
		downBytes, _ := templateFS.ReadFile("templates/migrations/migration_down.fizz")

		up = string(upBytes)
		down = string(downBytes)
	} else {
		migrationType = "sql"
	}

	err := boil.CreatePopMigration([]byte(up), []byte(down), name, migrationType)
	if err != nil {
		return err
	}

	return nil
}

func doMakeModel(modelName string) error {
	if modelName == "" {
		return errors.New("you must give the model a name")
	}

	data, err := templateFS.ReadFile("templates/data/model.go.txt")
	if err != nil {
		return err
	}

	model := string(data)

	plur := pluralize.NewClient()

	tableName := modelName

	if plur.IsPlural(modelName) {
		modelName = plur.Singular(modelName)
		tableName = strings.ToLower(tableName)
	} else {
		tableName = strings.ToLower(plur.Plural(modelName))
	}

	fileName := boil.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
	if fileExists(fileName) {
		return errors.New(fileName + " already exists!")
	}

	model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
	model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

	err = copyDataToFile([]byte(model), fileName)
	if err != nil {
		return err
	}

	return nil
}

func doMakeHandler(handlerName string) error {
	if handlerName == "" {
		return errors.New("you must give the handler a name")
	}

	fileName := boil.RootPath + "/handlers/" + strings.ToLower(handlerName) + ".go"
	if fileExists(fileName) {
		return errors.New(fileName + " already exists!")
	}

	data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
	if err != nil {
		return err
	}

	handler := string(data)
	handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(handlerName))

	err = ioutil.WriteFile(fileName, []byte(handler), 0o644)
	if err != nil {
		return err
	}

	return nil
}

func doMakeMail(mailName string) error {
	if mailName == "" {
		return errors.New("you must give the mail template a name")
	}

	htmlMail := boil.RootPath + "/mail/" + strings.ToLower(mailName) + ".html.tmpl"
	plainMail := boil.RootPath + "/mail/" + strings.ToLower(mailName) + ".plain.tmpl"

	err := copyFilefromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
	if err != nil {
		return err
	}

	err = copyFilefromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
	if err != nil {
		return err
	}

	return nil
}

func doSessionTable() error {
	dbType := boil.DB.DataType

	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	fileName := fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())

	upFile := boil.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := boil.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFilefromTemplate("templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table sessions"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	err = doMigrate("up", false)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}
