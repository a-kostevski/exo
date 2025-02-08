package fs

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// OSFileSystem implements the FileSystem interface using the os package.
type OSFileSystem struct{}

// NewOSFileSystem creates a new instance of OSFileSystem.
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// EnsureDirectoryExists ensures that the parent directory for the given file path exists.
// If you are dealing with a directory path, call EnsureDir directly.
func (fsys *OSFileSystem) EnsureDirectoryExists(path string) error {
	// We assume path is a file path; ensure its parent directory exists.
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0755)
}

// WriteFile writes content to the specified file. It ensures the parent directory exists.
func (fsys *OSFileSystem) WriteFile(path string, content []byte) error {
	if err := fsys.EnsureDirectoryExists(path); err != nil {
		return fmt.Errorf("failed to ensure directory exists for %s: %w", path, err)
	}
	return os.WriteFile(path, content, 0644)
}

// ReadFile reads and returns the contents of the specified file.
func (fsys *OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// FileExists returns true if the file at the given path exists.
func (fsys *OSFileSystem) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DeleteFile removes the file at the given path.
func (fsys *OSFileSystem) DeleteFile(path string) error {
	return os.Remove(path)
}

// AppendToFile appends the given content (with a newline) to the file at the specified path.
// If the file does not exist, it is created.
func (fsys *OSFileSystem) AppendToFile(path, content string) error {
	if err := fsys.EnsureDirectoryExists(path); err != nil {
		return fmt.Errorf("failed to ensure directory exists for %s: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer f.Close()
	_, err = io.WriteString(f, content+"\n")
	if err != nil {
		return fmt.Errorf("failed to append to file %s: %w", path, err)
	}
	return nil
}

// OpenInEditor opens the specified file in the given editor.
// It pipes the standard input/output and error streams to the editor process.
func (fsys *OSFileSystem) OpenInEditor(path, editor string) error {
	if path == "" {
		return fmt.Errorf("filepath cannot be empty")
	}
	if editor == "" {
		return fmt.Errorf("editor cannot be empty")
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (fsys *OSFileSystem) ReadDir(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}
