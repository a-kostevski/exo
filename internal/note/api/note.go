/*
Package api defines the core interfaces for the notes system. The package is designed
around the Interface Segregation Principle, breaking down note functionality into
smaller, focused interfaces that can be composed together.

The note system is built with the following design principles:

  - Interface Segregation: Core functionality is split into focused interfaces
  - Composition: Interfaces can be combined to create more complex behaviors
  - Extension: New note types can easily implement these interfaces
  - Separation of Concerns: Each interface handles a specific aspect of note functionality

The main interfaces are:

  - CoreNote: Basic note information (ID, title, content)
  - TemplateNote: Template-based content generation
  - ContentNote: Content manipulation
  - MetadataNote: Metadata management
  - PersistentNote: Storage operations

A complete note implementation should implement the Note interface, which combines
all the above interfaces. The system follows a factory pattern for note creation,
allowing easy registration and instantiation of different note types.

Example usage:

	// Create a new note through the factory
	factory := note.NewNoteFactory(templateManager)
	note, err := factory.Create(note.ZettelType, "My Note")
	if err != nil {
	    log.Fatal(err)
	}

	// Use the note
	note.SetContent("Hello, World!")
	if err := note.Save(); err != nil {
	    log.Fatal(err)
	}

The package is designed to be extended with new note types by implementing these
interfaces and registering them with the factory.
*/
package api

import "time"

// Note defines the core note interface
type Note interface {
	// Core information
	ID() string
	Title() string
	Content() string
	Path() string

	// Template operations
	ApplyTemplate(data interface{}) error

	// Content management
	SetContent(content string) error

	// Persistence
	Save() error
	Load() error
	Delete() error
	Open() error
}

// PeriodNote extends Note with time-based operations
type PeriodNote interface {
	Note
	PeriodStart() time.Time
	PeriodEnd() time.Time
	Previous() time.Time
	Next() time.Time
}

// NoteType identifies the type of note
type NoteType string

// PeriodType identifies the type of period note
type PeriodType string

const (
	// Note types
	ZettelType  NoteType = "zettel"
	ProjectType NoteType = "project"

	// Period types
	Daily PeriodType = "daily"
)

// NoteFactory defines the interface for creating notes
type NoteFactory interface {
	// Create creates a new note of the specified type
	Create(nType NoteType, title string) (Note, error)
}

// PeriodNoteFactory defines the interface for creating period notes
type PeriodNoteFactory interface {
	// Create creates a new period note of the specified type
	Create(pType PeriodType, date time.Time) (PeriodNote, error)
}
