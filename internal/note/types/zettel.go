/*
Package types provides specific implementations of different note types, building
on top of the base note implementation.

Each note type in this package specializes the base functionality for specific
use cases. The Zettelkasten implementation (ZettelNote) adds features specific
to the Zettelkasten method of note-taking, such as:

  - Reference management between notes
  - Tagging system
  - Specific validation rules
  - Metadata enrichment

The package demonstrates how to extend the base note implementation while
maintaining a clean separation of concerns. Each note type should:

  - Embed the base.BaseNote type
  - Implement additional methods specific to its use case
  - Override base methods where necessary
  - Provide proper validation
  - Handle its own metadata

Example usage:

	// Create a new Zettel note
	zet, err := types.NewZettelNote("My Thought", templateManager)
	if err != nil {
	    log.Fatal(err)
	}

	// Add Zettelkasten-specific metadata
	zet.AddTag("philosophy")
	zet.AddReference("20210101120000")

	// Save with enriched metadata
	if err := zet.Save(); err != nil {
	    log.Fatal(err)
	}
*/
package types

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/core"
	"github.com/a-kostevski/exo/internal/note/base"
	"github.com/a-kostevski/exo/internal/templates"
)

// ZettelNote represents a Zettelkasten note
type ZettelNote struct {
	*base.BaseNote
}

// NewZettelNote creates a new Zettelkasten note
func NewZettelNote(title string, tm templates.TemplateManager) (*ZettelNote, error) {
	baseNote := base.NewBaseNote(title, tm)

	if err := baseNote.SetTemplateName("zettel"); err != nil {
		return nil, fmt.Errorf("failed to set template: %w", err)
	}

	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	// Generate a sanitized file name from the title
	fileName := core.SanitizeFileName(title) + ".md"
	path := filepath.Join(cfg.ZettelDir, fileName)
	if err := baseNote.SetPath(path); err != nil {
		return nil, fmt.Errorf("failed to set path: %w", err)
	}

	note := &ZettelNote{
		BaseNote: baseNote,
	}

	// Apply template with data
	templateData := map[string]interface{}{
		"Title":    title,
		"Created":  time.Now().Format(time.RFC3339),
		"FileName": fileName,
	}

	if err := note.ApplyTemplate(templateData); err != nil {
		return nil, fmt.Errorf("failed to apply template: %w", err)
	}

	return note, nil
}
