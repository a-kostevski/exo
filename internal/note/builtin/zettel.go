package builtin

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/note"

	"github.com/a-kostevski/exo/internal/templates"
	"github.com/a-kostevski/exo/internal/utils"
)

// ZettelNote represents a Zettelkasten note
type ZettelNote struct {
	*note.BaseNote
}

// NewZettelNote creates a new Zettelkasten note
func NewZettelNote(title string, tm templates.TemplateManager) (*ZettelNote, error) {
	baseNote, err := note.NewBaseNote(title, tm)
	if err != nil {
		return nil, fmt.Errorf("failed to create base note: %w", err)
	}

	if err := baseNote.SetTemplateName("zettel"); err != nil {
		return nil, fmt.Errorf("failed to set template: %w", err)
	}

	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	// Generate a sanitized file name from the title
	fileName := utils.SanitizeFileName(title) + ".md"
	path := filepath.Join(cfg.ZettelDir, fileName)
	if err := baseNote.SetPath(path); err != nil {
		return nil, fmt.Errorf("failed to set path: %w", err)
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

func (n *ZettelNote) SetTitle(title string) error {
	return n.BaseNote.SetTitle(utils.SanitizeFileName(title))
}
