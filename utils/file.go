package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func CreateTempFile(initialContent string, extension string) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "restman_*."+extension)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	if _, err := tmpFile.Write([]byte(initialContent)); err != nil {
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	return tmpFile, nil
}

func RemoveTempFile(file *os.File) error {
	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	if err := os.Remove(file.Name()); err != nil {
		return fmt.Errorf("failed to remove temp file: %w", err)
	}
	return nil
}

func OpenInEditorCommand(file *os.File) *exec.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	return exec.Command(editor, file.Name())
}

func DownloadToTempFile(url string) (string, error) {
	// Fetch the schema from the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "openapi-*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Write the response body to the temp file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	// Return the temp file path
	return tmpFile.Name(), nil
}
