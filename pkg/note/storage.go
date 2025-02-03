package note

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
)

// Save saves the note content to disk
func (n *BaseNote) Save() error {
	logger.Debug("Saving note to disk",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path})

	if n.path == "" {
		logger.Error("Failed to save note: path not set",
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title})
		return errors.New("note path not set")
	}

	// Ensure directory exists
	dir := filepath.Dir(n.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create directory",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "path", Value: dir})
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write content to file
	if err := os.WriteFile(n.path, []byte(n.content), 0644); err != nil {
		logger.Error("Failed to write file",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "path", Value: n.path})
		return fmt.Errorf("failed to write file %s: %w", n.path, err)
	}

	logger.Info("Successfully saved note",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path})
	return nil
}

// Load reads the note content from disk
func (n *BaseNote) Load() error {
	logger.Debug("Loading note from disk",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path})

	if n.path == "" {
		logger.Error("Failed to load note: path not set",
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title})
		return errors.New("note path not set")
	}

	content, err := os.ReadFile(n.path)
	if err != nil {
		logger.Error("Failed to read file",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "path", Value: n.path})
		return fmt.Errorf("failed to read file %s: %w", n.path, err)
	}

	n.content = string(content)
	logger.Info("Successfully loaded note",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path})
	return nil
}

// Delete removes the note file from disk
func (n *BaseNote) Delete() error {
	logger.Debug("Deleting note",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path})

	if n.path == "" {
		logger.Error("Failed to delete note: path not set",
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title})
		return errors.New("note path not set")
	}

	if err := os.Remove(n.path); err != nil {
		logger.Error("Failed to delete file",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "path", Value: n.path})
		return fmt.Errorf("failed to delete file %s: %w", n.path, err)
	}

	logger.Info("Successfully deleted note",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path})
	return nil
}

// Exists checks if the note file exists on disk
func (n *BaseNote) Exists() bool {
	if n.path == "" {
		return false
	}
	_, err := os.Stat(n.path)
	return err == nil
}

// Open opens the note in the configured editor
func (n *BaseNote) Open() error {
	logger.Debug("Opening note in editor",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path},
	)
	if n.path == "" {
		logger.Error("Failed to open note: path not set",
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title})
		return errors.New("note path not set")
	}

	if !n.Exists() {
		logger.Error("Failed to open note: file does not exist",
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title},
			logger.Field{Key: "path", Value: n.path})
		return fmt.Errorf("note file does not exist: %s", n.path)
	}
	cfg := config.MustGet()
	if err := fs.OpenInEditor(n.path, cfg.General.Editor); err != nil {
		logger.Error("Failed to open note in editor",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title},
			logger.Field{Key: "path", Value: n.path})
		return fmt.Errorf("failed to open note in editor: %w", err)
	}

	logger.Info("Successfully opened note in editor",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "path", Value: n.path},
	)
	return nil
}
