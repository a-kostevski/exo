package note

import (
	"fmt"
	"time"

	"github.com/a-kostevski/exo/pkg/utils"
)

// ZettelNote represents a Zettelkasten note
type ZettelNote struct {
	*BaseNote
}

// NewZetNote creates a new Zettelkasten note
func NewZetNote(title string, opts ...NoteOption) (*ZettelNote, error) {
	fileName := utils.SanitizeFileName(title) + ".md"

	baseNote, err := NewBaseNote(title,
		WithFileName(fileName),
		WithSubDir("0-inbox"),
		WithTemplateName("zet"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create base note: %w", err)
	}

	note := &ZettelNote{
		BaseNote: baseNote,
	}

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
