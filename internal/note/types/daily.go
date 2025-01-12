package types

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/fs"
	"github.com/a-kostevski/exo/internal/note/api"
	"github.com/a-kostevski/exo/internal/note/base"
	"github.com/a-kostevski/exo/internal/templates"
)

// DailyNote represents a daily journal note
type DailyNote struct {
	*base.BasePeriodNote
}

// NewDailyNote creates a new daily note for the given date
func NewDailyNote(date time.Time, tm templates.TemplateManager) (*DailyNote, error) {
	title := date.Format("2006-01-02")
	baseNote := base.NewBasePeriodNote(title, api.Daily, date, tm)

	if err := baseNote.SetTemplateName("day"); err != nil {
		return nil, fmt.Errorf("failed to set template: %w", err)
	}

	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	path := filepath.Join(
		cfg.PeriodicDir,
		"day",
		baseNote.FormatPath(date),
	)

	if err := baseNote.SetPath(path); err != nil {
		return nil, fmt.Errorf("failed to set path: %w", err)
	}

	note := &DailyNote{
		BasePeriodNote: baseNote,
	}

	// Only apply template if the file doesn't exist
	if !fs.FileExists(path) {
		// Apply template with data
		templateData := map[string]interface{}{
			"Date":     date,
			"Previous": date.AddDate(0, 0, -1).Format("2006-01-02"),
			"Next":     date.AddDate(0, 0, 1).Format("2006-01-02"),
		}

		if err := note.ApplyTemplate(templateData); err != nil {
			return nil, fmt.Errorf("failed to apply template: %w", err)
		}
	} else {
		// Load existing content
		if err := note.Load(); err != nil {
			return nil, fmt.Errorf("failed to load existing note: %w", err)
		}
	}

	return note, nil
}
