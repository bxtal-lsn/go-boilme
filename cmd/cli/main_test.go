package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)

	// Execute the version command
	rootCmd.SetArgs([]string{"version"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check the output contains version info
	output := buf.String()
	if !strings.Contains(output, version) {
		t.Errorf("Expected output to contain version %s, got %s", version, output)
	}
}
