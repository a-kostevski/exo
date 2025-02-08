package zettel

import (
	"fmt"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/note"
	"github.com/a-kostevski/exo/pkg/templates"
)

// ZettelNote represents a specialized note (commonly known as a Zettel)
// that extends the basic functionality provided by BaseNote. In addition to
// the common note fields, a Zettel note includes a custom Tag field.
type ZettelNote struct {
	*note.BaseNote // Embed the core note functionality.
}

// NewZettelNote creates a new Zettel note with the specified title and tag.
// Dependencies are passed in (config, template manager, logger, fs) so that the
// note does not depend on global state. Default options (such as saving the note
// in the "zettel" subdirectory, using a filename based on the title, and applying
// the "zettel" template) are set; additional note options may be provided to
// override these defaults.
func NewZettelNote(title string, cfg config.Config, tm templates.TemplateManager, log logger.Logger, fs fs.FileSystem, opts ...note.NoteOption) (note.Note, error) {
	// Set defaults specific to Zettel notes.
	defaultOpts := []note.NoteOption{
		note.WithSubDir("0-inbox"),
		// For a default filename, we use the title with a ".md" extension.
		note.WithFileName(fmt.Sprintf("%s.md", title)),
		note.WithTemplateName("zet"),
	}
	// Merge the defaults with any options passed in.
	allOpts := append(defaultOpts, opts...)

	// Create the underlying BaseNote.
	base, err := note.NewBaseNote(title, cfg, tm, log, fs, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create base note: %w", err)
	}

	// Create the Zettel note by embedding the base note and setting the tag.
	zettel := &ZettelNote{
		BaseNote: base.(*note.BaseNote),
	}

	return zettel, nil
}

// Validate overrides the BaseNote's Validate method to enforce Zettel-specific rules.
// For example, it ensures that a tag is provided.
func (z *ZettelNote) Validate() error {
	if err := z.BaseNote.Validate(); err != nil {
		return err
	}
	return nil
}

// String returns a string representation of the Zettel note.
func (z *ZettelNote) String() string {
	return fmt.Sprintf("ZettelNote{Title: %s}", z.Title())
}

// Example method to update content and log the update.
func (z *ZettelNote) UpdateContent(newContent string) error {
	if err := z.SetContent(newContent); err != nil {
		return err
	}
	z.Logger.Infof("Updated content of Zettel note %s", z.Title())
	return nil
}

// Optionally, you might want to add custom behavior when saving a Zettel note.
// For example, you could override Save() if you need to perform additional steps.
// In this example, we simply use the BaseNote's Save() method.
func (z *ZettelNote) Save() error {
	// Here you could add custom pre-save logic.
	z.Logger.Infof("Saving Zettel note %s", z.Title())
	return z.BaseNote.Save()
}
