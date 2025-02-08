package note

import (
	"fmt"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
)

// NoteType distinguishes note types.
type NoteType string

// NoteFactory is an interface for creating a note.
type NoteFactory interface {
	CreateNote(title string, opts ...NoteOption) (Note, error)
	NoteType() NoteType
}

// BaseNoteFactory is a generic factory for BaseNotes.
type BaseNoteFactory struct {
	noteType NoteType
	Config   config.Config
	TM       templates.TemplateManager
	Logger   logger.Logger
	FS       fs.FileSystem
}

// NewBaseNoteFactory creates a new factory for a given note type.
func NewBaseNoteFactory(noteType NoteType, cfg config.Config, tm templates.TemplateManager, logger logger.Logger, fs fs.FileSystem) *BaseNoteFactory {
	return &BaseNoteFactory{
		noteType: noteType,
		Config:   cfg,
		TM:       tm,
		Logger:   logger,
		FS:       fs,
	}
}

func (f *BaseNoteFactory) NoteType() NoteType {
	return f.noteType
}

func (f *BaseNoteFactory) CreateNote(title string, opts ...NoteOption) (Note, error) {
	note, err := NewBaseNote(title, f.Config, f.TM, f.Logger, f.FS, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create base note: %w", err)
	}
	return note, nil
}
