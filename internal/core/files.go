/*
Package core provides fundamental operations and services that the application relies on.
It includes application-specific functionality that integrates with the configuration
and other application components.

The package focuses on providing:
  - Editor integration with configuration awareness
  - File name sanitization for note creation
  - Application-specific file operations
  - Integration with other application components

These operations are considered "core" because they:
  - Have application-specific requirements
  - Need to be consistent across the application
  - Integrate with the application's configuration

Example usage:

	// Create a note with sanitized filename
	name := core.SanitizeFileName("My Note Title")
	path := filepath.Join(config.Get().NotesDir, name + ".md")

	// Open in configured editor
	if err := core.OpenInEditor(path); err != nil {
		return err
	}
*/
package core

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/a-kostevski/exo/internal/config"
)

// OpenInEditor opens a file in the configured editor.
// It uses the editor specified in the application configuration.
// If no editor is configured, it returns without error.
func OpenInEditor(filepath string, cfg *config.Config) error {
	if filepath == "" {
		return fmt.Errorf("filepath cannot be empty")
	}

	editor := cfg.Editor
	if editor == "" {
		return nil
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open in editor: %w", err)
	}

	return nil
}

// SanitizeFileName converts a string into a safe filename by:
// - Converting spaces to dashes
// - Removing non-alphanumeric characters (except dashes)
// - Converting to lowercase
func SanitizeFileName(name string) string {
	// Replace spaces with dashes
	name = strings.ReplaceAll(name, " ", "-")

	// Remove non-alphanumeric characters (except dashes)
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-]`)
	name = reg.ReplaceAllString(name, "")

	return strings.ToLower(name)
}
