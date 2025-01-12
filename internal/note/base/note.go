/*
Package base provides foundational implementations of the note system's interfaces.
It implements the core functionality that can be reused across different note types.

The BaseNote type provides a complete implementation of the api.Note interface,
handling common operations such as:

  - Basic note information management (ID, title, content)
  - Template-based content generation
  - Metadata storage and retrieval
  - File-based persistence
  - Modification tracking

BaseNote is designed to be embedded in more specific note implementations, following
the composition pattern. This allows specialized note types to inherit basic
functionality while adding their own specific features.

Example usage:

	// Create a specialized note type
	type MyNote struct {
	    *base.BaseNote
	    specialField string
	}

	// Create a new instance
	myNote := &MyNote{
	    BaseNote: base.NewBaseNote("title", templateManager),
	    specialField: "value",
	}

The base package focuses on providing robust, well-tested implementations of common
note operations, allowing specialized note types to focus on their unique features
rather than reimplementing basic functionality.
*/
package base

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/core"
	"github.com/a-kostevski/exo/internal/fs"
	"github.com/a-kostevski/exo/internal/templates"
)

// BaseNote provides a basic implementation of the Note interface
type BaseNote struct {
	id           string
	title        string
	content      string
	path         string
	templateName string
	created      time.Time
	modified     time.Time
	tm           templates.TemplateManager
}

// NewBaseNote creates a new BaseNote with the given title
func NewBaseNote(title string, tm templates.TemplateManager) *BaseNote {
	now := time.Now()
	return &BaseNote{
		id:       fmt.Sprintf("%d", now.UnixNano()),
		title:    title,
		created:  now,
		modified: now,
		tm:       tm,
	}
}

// Core note information
func (n *BaseNote) ID() string      { return n.id }
func (n *BaseNote) Title() string   { return n.title }
func (n *BaseNote) Content() string { return n.content }
func (n *BaseNote) Path() string    { return n.path }

// Template operations
func (n *BaseNote) ApplyTemplate(data interface{}) error {
	if n.templateName == "" {
		return errors.New("no template name set")
	}
	if n.tm == nil {
		return errors.New("template manager not initialized")
	}

	content, err := n.tm.ProcessTemplate(n.templateName, data)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	n.content = content
	n.modified = time.Now()
	return nil
}

// Content management
func (n *BaseNote) SetContent(content string) error {
	n.content = content
	n.modified = time.Now()
	return nil
}

// SetPath sets the note's path
func (n *BaseNote) SetPath(path string) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}
	n.path = path
	return nil
}

// SetTemplateName sets the template to use
func (n *BaseNote) SetTemplateName(name string) error {
	if name == "" {
		return errors.New("template name cannot be empty")
	}
	n.templateName = name
	return nil
}

// Persistence
func (n *BaseNote) Save() error {
	if n.path == "" {
		return errors.New("note path not set")
	}

	fs.EnsureDirectoryExists(filepath.Dir(n.path))
	return fs.WriteFile(n.path, []byte(n.content))
}

func (n *BaseNote) Load() error {
	if n.path == "" {
		return errors.New("note path not set")
	}

	content, err := fs.ReadFile(n.path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	n.content = string(content)
	return nil
}

func (n *BaseNote) Delete() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	return fs.DeleteFile(n.path)
}

// Open opens the note in the configured editor
func (n *BaseNote) Open() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	cfg := config.MustGet()
	return core.OpenInEditor(n.path, cfg)
}
