package cmd

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Generate authentication system",
	Long: `Generate a complete authentication system including:
  - Database tables
  - Models
  - Handlers
  - Middleware
  - Views`,
	Run: func(cmd *cobra.Command, args []string) {
		// Ask for confirmation
		var confirmAuth bool
		prompt := &survey.Confirm{
			Message: "This will create authentication tables, models, and handlers. Continue?",
			Default: true,
		}
		survey.AskOne(prompt, &confirmAuth)

		if !confirmAuth {
			color.Yellow("Auth setup cancelled")
			return
		}

		// Start spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Setting up authentication system..."
		s.Color("green")
		s.Start()

		// Setup necessary for database operations
		setupRoot(cmd, args)

		err := doAuth()

		s.Stop()
		if err != nil {
			exitGracefully(err)
		}

		color.Green("✓ Authentication system created successfully!")
		fmt.Println("\nThe following components were created:")
		color.Cyan("  - Database tables: users, tokens, remember_tokens")
		color.Cyan("  - Models: User, Token, RememberToken")
		color.Cyan("  - Middleware: Auth, AuthToken, Remember")
		color.Cyan("  - Handlers: Login, Logout, Forgot Password, Reset Password")
		color.Cyan("  - Views: Login, Forgot Password, Reset Password")

		color.Yellow("\nNext steps:")
		fmt.Println("1. Add the models to your data/models.go file")
		fmt.Println("2. Add the middleware to your routes")
		fmt.Println("3. Add the routes to your routes.go file")
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func doAuth() error {
	// Your existing auth functionality from auth.go
	checkForDB()

	// migrations
	dbType := boil.DB.DataType

	tx, err := boil.PopConnect()
	if err != nil {
		return err
	}
	defer tx.Close()

	upBytes, err := templateFS.ReadFile("templates/migrations/auth_tables." + dbType + ".sql")
	if err != nil {
		return err
	}

	downBytes := []byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens;")

	err = boil.CreatePopMigration(upBytes, downBytes, "auth", "sql")
	if err != nil {
		return err
	}

	// run migrations
	err = boil.RunPopMigrations(tx)
	if err != nil {
		return err
	}

	// Copy files
	filesToCopy := []struct {
		src string
		dst string
	}{
		{"templates/data/user.go.txt", boil.RootPath + "/data/user.go"},
		{"templates/data/token.go.txt", boil.RootPath + "/data/token.go"},
		{"templates/data/remember_token.go.txt", boil.RootPath + "/data/remember_token.go"},
		{"templates/middleware/auth.go.txt", boil.RootPath + "/middleware/auth.go"},
		{"templates/middleware/auth-token.go.txt", boil.RootPath + "/middleware/auth-token.go"},
		{"templates/middleware/remember.go.txt", boil.RootPath + "/middleware/remember.go"},
		{"templates/handlers/auth-handlers.go.txt", boil.RootPath + "/handlers/auth-handlers.go"},
		{"templates/mailer/password-reset.html.tmpl", boil.RootPath + "/mail/password-reset.html.tmpl"},
		{"templates/mailer/password-reset.plain.tmpl", boil.RootPath + "/mail/password-reset.plain.tmpl"},
		{"templates/views/login.jet", boil.RootPath + "/views/login.jet"},
		{"templates/views/forgot.jet", boil.RootPath + "/views/forgot.jet"},
		{"templates/views/reset-password.jet", boil.RootPath + "/views/reset-password.jet"},
	}

	for _, file := range filesToCopy {
		err = copyFilefromTemplate(file.src, file.dst)
		if err != nil {
			return err
		}
	}

	return nil
}
