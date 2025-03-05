package cmd

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
)

//go:embed templates
var templateFS embed.FS

// checkForDB checks if a database connection is available
func checkForDB() {
	dbType := boil.DB.DataType

	if dbType == "" {
		exitGracefully(errors.New("no database connection provided in .env"))
	}

	if !fileExists(boil.RootPath + "/config/database.yml") {
		exitGracefully(errors.New("config/database.yml does not exist"))
	}
}

// copyFilefromTemplate copies a file from the embedded templates to the target location
func copyFilefromTemplate(templatePath, targetFile string) error {
	if fileExists(targetFile) {
		return errors.New(targetFile + " already exists!")
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

// copyDataToFile copies byte data to a file
func copyDataToFile(data []byte, to string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(to)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	err := os.WriteFile(to, data, 0o644)
	if err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file exists
func fileExists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}
