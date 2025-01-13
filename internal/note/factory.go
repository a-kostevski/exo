package note

import (
	"fmt"
	"time"

	"github.com/a-kostevski/exo/internal/templates"
)

// Factory is responsible for creating instances of different note types.
type Factory struct {
	tm           templates.TemplateManager
	constructors map[string]Constructor
}

// Constructor is a function type for creating notes
type Constructor func(title string, date time.Time, tm templates.TemplateManager) (Note, error)

// NewFactory creates a new Factory instance.
func NewFactory(tm templates.TemplateManager) *Factory {
	return &Factory{
		tm:           tm,
		constructors: make(map[string]Constructor),
	}
}

// Register adds a note type constructor to the factory
func (f *Factory) Register(noteType string, constructor Constructor) {
	f.constructors[noteType] = constructor
}

// New creates a new note of the specified type
func (f *Factory) New(noteType, title string, date time.Time) (Note, error) {
	constructor, ok := f.constructors[noteType]
	if !ok {
		return nil, fmt.Errorf("unsupported note type: %s", noteType)
	}
	return constructor(title, date, f.tm)
}
