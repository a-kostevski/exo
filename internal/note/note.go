package note

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/templates"
	"github.com/a-kostevski/exo/internal/utils"
)

type Note interface {
	ID() string
	Title() string
	Content() string
	Path() string

	SetContent(content string) error
	SetTitle(title string) error

	Save() error
	Load() error
	Delete() error

	ApplyTemplate(data interface{}) error
}

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

func NewBaseNote(title string, tm templates.TemplateManager) (*BaseNote, error) {
	if tm == nil {
		return nil, errors.New("template manager not provided")
	}
	return &BaseNote{
		id:       fmt.Sprintf("%d", time.Now().UnixNano()),
		title:    title,
		created:  time.Now(),
		modified: time.Now(),
		tm:       tm,
	}, nil
}

// Core note information
func (n *BaseNote) ID() string      { return n.id }
func (n *BaseNote) Title() string   { return n.title }
func (n *BaseNote) Content() string { return n.content }
func (n *BaseNote) Path() string    { return n.path }

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
	return n.SetContent(content)
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

	if err := utils.EnsureDirectoryExists(filepath.Dir(n.path)); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	return utils.WriteFile(n.path, []byte(n.content))
}

func (n *BaseNote) Load() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	content, err := utils.ReadFile(n.path)
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
	return utils.DeleteFile(n.path)
}

// Open opens the note in the configured editor
func (n *BaseNote) Open() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	cfg := config.MustGet()
	return utils.OpenInEditor(n.path, cfg.Editor)
}

func (n *BaseNote) SetTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	n.title = title
	n.modified = time.Now()
	return nil
}
