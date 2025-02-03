package note

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
)

type BaseNote struct {
	id           string
	title        string
	content      string
	path         string
	fileName     string
	subDir       string
	templateName string
	created      time.Time
	modified     time.Time
}

type BaseNoteOpts struct {
	Title        string
	Content      string
	TemplateName string
	SubDir       string
	FileName     string
}

type NoteOption func(*BaseNote) error

func NewBaseNote(title string, opts ...NoteOption) (*BaseNote, error) {
	logger.Debug("Creating new base note",
		logger.Field{Key: "title", Value: title})

	if title == "" {
		logger.Error("Failed to create note: empty title")
		return nil, errors.New("title cannot be empty")
	}

	note := &BaseNote{
		title:    title,
		created:  time.Now(),
		modified: time.Now(),
	}

	for _, opt := range opts {
		if err := opt(note); err != nil {
			logger.Error("Failed to apply note option",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "title", Value: title})
			return nil, fmt.Errorf("failed to create note: %w", err)
		}
	}

	logger.Info("Base note created successfully",
		logger.Field{Key: "id", Value: note.id},
		logger.Field{Key: "title", Value: note.title},
		logger.Field{Key: "created", Value: note.created})

	return note, nil
}

// Core note information
func (n *BaseNote) ID() string          { return n.id }
func (n *BaseNote) Title() string       { return n.title }
func (n *BaseNote) Content() string     { return n.content }
func (n *BaseNote) Path() string        { return n.path }
func (n *BaseNote) Created() time.Time  { return n.created }
func (n *BaseNote) Modified() time.Time { return n.modified }

// Content management
func (n *BaseNote) SetContent(content string) error {
	logger.Debug("Setting note content",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title})

	n.content = content
	n.modified = time.Now()

	logger.Debug("Note content updated",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "modified", Value: n.modified})
	return nil
}

// Template operations
func (n *BaseNote) ApplyTemplate(data interface{}) error {
	logger.Debug("Applying template to note",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "template", Value: n.templateName})

	if n.templateName == "" {
		logger.Error("Failed to apply template: no template name set",
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title})
		return errors.New("no template name set")
	}

	cfg := config.MustGet()
	tm, err := templates.NewTemplateManager(templates.TemplateConfig{
		TemplateDir:       cfg.Dir.TemplateDir,
		TemplateExtension: ".md",
		Logger:            logger.Default(), // Add a way to get default logger
		FilePermissions:   0644,
	})
	if err != nil {
		logger.Error("Failed to create template manager",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title})
		return fmt.Errorf("failed to create template manager: %w", err)
	}

	content, err := tm.ProcessTemplate(n.templateName, data)
	if err != nil {
		logger.Error("Failed to process template",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "id", Value: n.id},
			logger.Field{Key: "title", Value: n.title},
			logger.Field{Key: "template", Value: n.templateName})
		return fmt.Errorf("failed to process template: %w", err)
	}

	logger.Info("Successfully applied template to note",
		logger.Field{Key: "id", Value: n.id},
		logger.Field{Key: "title", Value: n.title},
		logger.Field{Key: "template", Value: n.templateName})
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

func (n *BaseNote) SetTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	n.title = title
	n.modified = time.Now()
	return nil
}

func (n *BaseNote) String() string {
	return fmt.Sprintf("Note{ID: %s, Title: %s}", n.id, n.title)
}

// Add validation for required fields
func (n *BaseNote) Validate() error {
	if n.title == "" {
		return errors.New("title is required")
	}
	if n.path == "" {
		return errors.New("path is required")
	}
	return nil
}

// WithSubDir sets the subdirectory for the note relative to data home
func WithSubDir(subDir string) NoteOption {
	return func(n *BaseNote) error {
		if subDir == "" {
			return errors.New("subdirectory cannot be empty")
		}
		n.subDir = subDir
		return n.updatePath()
	}
}

// WithFileName sets the filename for the note
func WithFileName(fileName string) NoteOption {
	return func(n *BaseNote) error {
		if fileName == "" {
			return errors.New("filename cannot be empty")
		}
		n.fileName = fileName
		return n.updatePath()
	}
}

// updatePath constructs the full path from DataHome, SubDir and FileName
func (n *BaseNote) updatePath() error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Construct path relative to DataHome
	n.path = filepath.Join(
		cfg.Dir.DataHome,
		n.subDir,
		n.fileName,
	)
	return nil
}

func WithTemplateName(name string) NoteOption {
	return func(n *BaseNote) error {
		return n.SetTemplateName(name)
	}
}

func WithContent(content string) NoteOption {
	return func(n *BaseNote) error {
		return n.SetContent(content)
	}
}
