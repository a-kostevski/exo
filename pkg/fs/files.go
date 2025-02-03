package fs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/a-kostevski/exo/pkg/logger"
)

// EnsureDirectoryExists creates directory if it doesn't exist.
// It creates all necessary parent directories with appropriate permissions (0755).
func EnsureDirectoryExists(path string) error {
	dir := filepath.Dir(path)
	logger.Debug("Ensuring directory exists", logger.Field{Key: "path", Value: dir})

	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create directory",
			logger.Field{Key: "path", Value: dir},
			logger.Field{Key: "error", Value: err})
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	logger.Debug("Directory created successfully", logger.Field{Key: "path", Value: dir})
	return nil
}

// FileExists verifies if a file exists at the given path.
// It returns true if the file exists and is accessible.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	exists := err == nil || os.IsExist(err)
	logger.Debug("Checking file existence",
		logger.Field{Key: "path", Value: path},
		logger.Field{Key: "exists", Value: exists})
	return exists
}

// CreateFile creates a new empty file at the specified path.
// It also ensures that the parent directory exists.
func CreateFile(path string) error {
	logger.Info("Creating new file", logger.Field{Key: "path", Value: path})

	dir := filepath.Dir(path)
	if err := EnsureDirectoryExists(dir); err != nil {
		logger.Error("Failed to create directory for file",
			logger.Field{Key: "path", Value: path},
			logger.Field{Key: "error", Value: err})
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		logger.Error("Failed to create file",
			logger.Field{Key: "path", Value: path},
			logger.Field{Key: "error", Value: err})
		return err
	}
	defer f.Close()

	logger.Info("File created successfully", logger.Field{Key: "path", Value: path})
	return nil
}

// WriteFile writes content to a file with appropriate permissions (0644).
// If the file exists, it will be overwritten.
func WriteFile(path string, content []byte) error {
	logger.Info("Writing to file",
		logger.Field{Key: "path", Value: path},
		logger.Field{Key: "size", Value: len(content)})

	if err := EnsureDirectoryExists(path); err != nil {
		logger.Error("Failed to ensure directory exists",
			logger.Field{Key: "path", Value: path},
			logger.Field{Key: "error", Value: err})
		return err
	}

	if err := os.WriteFile(path, content, 0644); err != nil {
		logger.Error("Failed to write file",
			logger.Field{Key: "path", Value: path},
			logger.Field{Key: "error", Value: err})
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	logger.Info("File written successfully", logger.Field{Key: "path", Value: path})
	return nil
}

// ReadFile reads the entire contents of a file.
// It returns the content as a byte slice.
func ReadFile(path string) ([]byte, error) {
	logger.Debug("Reading file", logger.Field{Key: "path", Value: path})

	content, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Failed to read file",
			logger.Field{Key: "path", Value: path},
			logger.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	logger.Debug("File read successfully",
		logger.Field{Key: "path", Value: path},
		logger.Field{Key: "size", Value: len(content)})
	return content, nil
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
	logger.Info("Deleting file", logger.Field{Key: "path", Value: path})

	if err := os.Remove(path); err != nil {
		if !os.IsNotExist(err) {
			logger.Error("Failed to delete file",
				logger.Field{Key: "path", Value: path},
				logger.Field{Key: "error", Value: err})
			return fmt.Errorf("failed to delete file %s: %w", path, err)
		}
		logger.Warn("File to delete does not exist",
			logger.Field{Key: "path", Value: path})
	}

	logger.Info("File deleted successfully", logger.Field{Key: "path", Value: path})
	return nil
}

// OpenInEditor opens a file in the configured editor.
// It uses the editor specified in the application configuration.
// If no editor is configured, it returns without error.
func OpenInEditor(filepath string, editor string) error {
	logger.Info("Opening file in editor",
		logger.Field{Key: "filepath", Value: filepath},
		logger.Field{Key: "editor", Value: editor})

	if filepath == "" {
		logger.Error("Empty filepath provided")
		return fmt.Errorf("filepath cannot be empty")
	}

	if editor == "" {
		logger.Error("Empty editor provided")
		return fmt.Errorf("editor cannot be empty")
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Debug("Executing editor command",
		logger.Field{Key: "command", Value: cmd.String()})

	if err := cmd.Run(); err != nil {
		logger.Error("Failed to open in editor",
			logger.Field{Key: "filepath", Value: filepath},
			logger.Field{Key: "editor", Value: editor},
			logger.Field{Key: "error", Value: err})
		return fmt.Errorf("failed to open in editor: %w", err)
	}

	logger.Info("File opened successfully in editor",
		logger.Field{Key: "filepath", Value: filepath},
		logger.Field{Key: "editor", Value: editor})
	return nil
}

// SanitizeFileName converts a string into a safe filename by:
// - Converting spaces to dashes
// - Removing non-alphanumeric characters (except dashes)
// - Converting to lowercase
func SanitizeFileName(name string) string {
	logger.Debug("Sanitizing filename",
		logger.Field{Key: "original", Value: name})

	sanitized := name
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-]`)
	sanitized = reg.ReplaceAllString(sanitized, "")
	sanitized = strings.ToLower(sanitized)

	logger.Debug("Filename sanitized",
		logger.Field{Key: "original", Value: name},
		logger.Field{Key: "sanitized", Value: sanitized})
	return sanitized
}

// EnsureDirectories creates all necessary directories with appropriate permissions.
// It expands and sanitizes paths using the user's home directory as base.
// It returns an error if any directory cannot be created.
func EnsureDirectories(paths ...string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	for _, path := range paths {
		absPath := SanitizePath(path, home)
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", absPath, err)
		}
		logger.Info("Created directory", logger.Field{Key: "path", Value: absPath})
	}
	return nil
}
