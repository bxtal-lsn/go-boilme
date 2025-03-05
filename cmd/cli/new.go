package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var appURL string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [app_name]",
	Short: "Create a new Boilme application",
	Long: `Create a new Boilme application with the specified name.
This will clone the starter template and set up everything you need.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		// Start spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Creating new Boilme application..."
		s.Color("green")
		s.Start()

		err := doNew(appName)

		s.Stop()
		if err != nil {
			exitGracefully(err)
		}

		color.Green("✓ Successfully created new Boilme application: %s", appName)
		color.Yellow("To get started:")
		fmt.Printf("  cd %s\n", appName)
		fmt.Println("  go mod tidy")
		fmt.Println("  go run .")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func doNew(appName string) error {
	appName = strings.ToLower(appName)
	appURL = appName

	// sanitize the application name (convert url to single word)
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[(len(exploded) - 1)]
	}

	// git clone the skeleton application
	_, err := git.PlainClone("./"+appName, false, &git.CloneOptions{
		URL:      "https://github.com/bxtal-lsn/go-boilme-app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		return err
	}

	// remove .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", appName))
	if err != nil {
		return err
	}

	// create a ready to go .env file
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		return err
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	env = strings.ReplaceAll(env, "${KEY}", boil.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		return err
	}

	// create a makefile
	if runtime.GOOS == "windows" {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.windows", appName))
		if err != nil {
			return err
		}
		defer source.Close()

		destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			return err
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}
	} else {
		source, err := os.Open(fmt.Sprintf("./%s/Makefile.mac", appName))
		if err != nil {
			return err
		}
		defer source.Close()

		destination, err := os.Create(fmt.Sprintf("./%s/Makefile", appName))
		if err != nil {
			return err
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
		if err != nil {
			return err
		}
	}
	_ = os.Remove("./" + appName + "/Makefile.mac")
	_ = os.Remove("./" + appName + "/Makefile.windows")

	// update the go.mod file
	_ = os.Remove("./" + appName + "/go.mod")

	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		return err
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(mod), "./"+appName+"/go.mod")
	if err != nil {
		return err
	}

	// update existing .go files with correct name/imports
	os.Chdir("./" + appName)
	updateSource()

	// run go mod tidy in the project directory
	cmd := exec.Command("go", "get", "github.com/bxtal-lsn/go-boilme")
	err = cmd.Start()
	if err != nil {
		return err
	}

	cmd = exec.Command("go", "mod", "tidy")
	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func updateSource() error {
	// Walk entire project folder, including subfolders
	err := filepath.Walk(".", func(path string, fi os.FileInfo, err error) error {
		// Check for errors before doing anything else
		if err != nil {
			return err
		}

		// Check if current file is directory
		if fi.IsDir() {
			return nil
		}

		// Only check go files
		matched, err := filepath.Match("*.go", fi.Name())
		if err != nil {
			return err
		}

		// We have a matching file
		if matched {
			// Read file contents
			read, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			newContents := strings.Replace(string(read), "myapp", appURL, -1)

			// Write the changed file
			err = os.WriteFile(path, []byte(newContents), 0o644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
