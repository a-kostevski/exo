package note

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
)

// Note defines the interface that all note types must satisfy.
type Note interface {
	// Content management
	Content() string
	SetContent(string) error

	// Filesystem operations
	Save() error
	Load() error
	Delete() error
	Exists() bool

	// Editor operation
	Open() error

	// Metadata accessors
	Title() string
	Path() string
	Created() time.Time
	Modified() time.Time

	// Validate the note (e.g., check that required fields are set)
	Validate() error

	// String returns a summary string.
	String() string
}

// BaseNote holds common fields and dependencies for a note.
type BaseNote struct {
	// Core attributes
	id           string
	title        string
	content      string
	path         string // full file path
	fileName     string
	subDir       string
	templateName string

	created  time.Time
	modified time.Time

	// Dependencies (injected via the constructor)
	Config config.Config
	TM     templates.TemplateManager
	Logger logger.Logger
	FS     fs.FileSystem
}

// NoteOption defines a functional option for configuring a BaseNote.
type NoteOption func(*BaseNote) error

// NewBaseNote creates a new BaseNote instance with the given title and dependencies.
// Additional options (like setting the subdirectory, filename, template, etc.) can be provided.
func NewBaseNote(title string, cfg config.Config, tm templates.TemplateManager, logger logger.Logger, fs fs.FileSystem, opts ...NoteOption) (Note, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}

	n := &BaseNote{
		title:    title,
		created:  time.Now(),
		modified: time.Now(),
		Config:   cfg,
		TM:       tm,
		Logger:   logger,
		FS:       fs,
	}

	// Apply functional options to set additional attributes.
	for _, opt := range opts {
		if err := opt(n); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	// updatePath must be called if both subDir and fileName are set.
	if n.subDir == "" || n.fileName == "" {
		return nil, errors.New("subdirectory and filename must be provided")
	}
	if err := n.updatePath(); err != nil {
		return nil, err
	}

	return n, nil
}

// WithSubDir sets the subdirectory for the note.
func WithSubDir(subDir string) NoteOption {
	return func(n *BaseNote) error {
		if subDir == "" {
			return errors.New("subdirectory cannot be empty")
		}
		n.subDir = subDir
		return nil
	}
}

// WithFileName sets the filename for the note.
func WithFileName(fileName string) NoteOption {
	return func(n *BaseNote) error {
		if fileName == "" {
			return errors.New("filename cannot be empty")
		}
		n.fileName = fileName
		return nil
	}
}

// WithTemplateName sets the name of the template to be applied.
func WithTemplateName(templateName string) NoteOption {
	return func(n *BaseNote) error {
		if templateName == "" {
			return errors.New("template name cannot be empty")
		}
		n.templateName = templateName
		return nil
	}
}

// WithContent sets initial content.
func WithContent(content string) NoteOption {
	return func(n *BaseNote) error {
		n.content = content
		return nil
	}
}

// updatePath calculates the full file path based on the configuration, subdirectory, and filename.
func (n *BaseNote) updatePath() error {
	n.path = filepath.Join(n.Config.Dir.DataHome, n.subDir, n.fileName)
	return nil
}

// Implement the Note interface:

func (n *BaseNote) Content() string {
	return n.content
}

func (n *BaseNote) SetContent(content string) error {
	n.content = content
	n.modified = time.Now()
	return nil
}

func (n *BaseNote) Save() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	// Ensure the parent directory exists.
	if err := n.FS.EnsureDirectoryExists(n.path); err != nil {
		return err
	}
	if err := os.WriteFile(n.path, []byte(n.content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", n.path, err)
	}
	return nil
}

func (n *BaseNote) Load() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	content, err := os.ReadFile(n.path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", n.path, err)
	}
	n.content = string(content)
	return nil
}

func (n *BaseNote) Delete() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	if err := os.Remove(n.path); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", n.path, err)
	}
	return nil
}

func (n *BaseNote) Exists() bool {
	return n.FS.FileExists(n.path)
}

func (n *BaseNote) Open() error {
	if n.path == "" {
		return errors.New("note path not set")
	}
	if !n.Exists() {
		return fmt.Errorf("note file does not exist: %s", n.path)
	}
	return n.FS.OpenInEditor(n.path, n.Config.General.Editor)
}

func (n *BaseNote) Title() string {
	return n.title
}

func (n *BaseNote) Path() string {
	return n.path
}

func (n *BaseNote) Created() time.Time {
	return n.created
}

func (n *BaseNote) Modified() time.Time {
	return n.modified
}

func (n *BaseNote) Validate() error {
	if n.title == "" {
		return errors.New("title is required")
	}
	if n.path == "" {
		return errors.New("path is required")
	}
	return nil
}

func (n *BaseNote) String() string {
	return fmt.Sprintf("Note{ID: %s, Title: %s}", n.id, n.title)
}

// ApplyTemplate uses the template manager to process a template and sets the note content.
func (n *BaseNote) ApplyTemplate(data interface{}) error {
	if n.templateName == "" {
		return errors.New("no template name set")
	}
	content, err := n.TM.ProcessTemplate(n.templateName, data)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}
	return n.SetContent(content)
}
