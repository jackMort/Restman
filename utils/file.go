package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenInEditor(initialContent string, extension string) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "restman_*."+extension)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name()) // Clean up temp file
	}()

	// Write the initial content to the temp file
	if _, err := tmpFile.Write([]byte(initialContent)); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	// Determine editor to use
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Open the file in the editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read the edited content back from the file using os.ReadFile
	editedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	return string(editedContent), nil
}
