package fs

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureDirectoryExists creates directory if it doesn't exist.
// It creates all necessary parent directories with appropriate permissions (0755).
func EnsureDirectoryExists(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}

// FileExists verifies if a file exists at the given path.
// It returns true if the file exists and is accessible.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// CreateFile creates a new empty file at the specified path.
// It also ensures that the parent directory exists.
func CreateFile(path string) error {
	dir := filepath.Dir(path)

	if err := EnsureDirectoryExists(dir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	return f.Close()
}

// WriteFile writes content to a file with appropriate permissions (0644).
// If the file exists, it will be overwritten.
func WriteFile(path string, content []byte) error {
	if err := EnsureDirectoryExists(path); err != nil {
		return err
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

// ReadFile reads the entire contents of a file.
// It returns the content as a byte slice.
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return data, nil
}

// AppendToFile appends content to an existing file.
// If the file doesn't exist, it will be created.
// A newline is automatically added after the content.
func AppendToFile(path, content string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString(content + "\n"); err != nil {
		return fmt.Errorf("failed to append to file: %w", err)
	}
	return nil
}

// DeleteFile removes a file if it exists
func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}
	return nil
}
