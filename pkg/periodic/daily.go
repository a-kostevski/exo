package periodic

import (
	"fmt"
	"time"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/note"
	"github.com/a-kostevski/exo/pkg/templates"
)

// DailyNavigator implements PeriodNavigator for daily notes.
type DailyNavigator struct{}

func (dn *DailyNavigator) Previous(date time.Time) time.Time {
	return date.AddDate(0, 0, -1)
}

func (dn *DailyNavigator) Next(date time.Time) time.Time {
	return date.AddDate(0, 0, 1)
}

func (dn *DailyNavigator) Start(date time.Time) time.Time {
	// For a daily note, you might simply use the date as-is.
	return date
}

func (dn *DailyNavigator) End(date time.Time) time.Time {
	// The "end" of the day can be defined as the same day (or at 23:59:59 if needed).
	return date
}

// DailyNote represents a daily periodic note.
type DailyNote struct {
	*PeriodicNote // Embeds all periodic note functionality.
}

// NewDailyNote creates (or loads) a daily note for the given date.
// It sets daily-specific defaults (e.g. subdirectory "day", filename based on date, template "day")
// and sets the navigator to a DailyNavigator.
func NewDailyNote(date time.Time, cfg config.Config, tm templates.TemplateManager, log logger.Logger, fs fs.FileSystem) (*DailyNote, error) {
	// For a daily note, use the date formatted as YYYY-MM-DD as the title.
	title := date.Format("2006-01-02")
	// Set defaults: place the note in a "day" subdirectory, use a file name "<date>.md",
	// and choose the "day" template.
	opts := []note.NoteOption{
		note.WithSubDir("day"),
		note.WithFileName(fmt.Sprintf("%s.md", title)),
		note.WithTemplateName("day"),
	}
	// Create the underlying PeriodicNote.
	p, err := NewPeriodicNote(title, date, cfg, tm, log, fs, opts...)
	if err != nil {
		log.Error("Failed to create periodic note",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "title", Value: title})
		return nil, fmt.Errorf("failed to create periodic note: %w", err)
	}
	// Set the navigator to DailyNavigator.
	p.SetNavigator(&DailyNavigator{})

	daily := &DailyNote{
		PeriodicNote: p,
	}

	// If the note file does not exist, initialize its content.
	if !daily.Exists() {
		log.Info("Initializing new daily note",
			logger.Field{Key: "path", Value: daily.Path()})
		templateData := map[string]interface{}{
			"Date":     title,
			"Previous": daily.PreviousOrZero().Format("2006-01-02"),
			"Next":     daily.NextOrZero().Format("2006-01-02"),
		}
		if err := daily.ApplyTemplate(templateData); err != nil {
			log.Error("Failed to apply template",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "path", Value: daily.Path()})
			return nil, fmt.Errorf("failed to apply template: %w", err)
		}
		if err := daily.Save(); err != nil {
			log.Error("Failed to save daily note",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "path", Value: daily.Path()})
			return nil, fmt.Errorf("failed to save daily note: %w", err)
		}
	} else {
		// Otherwise, load the existing note.
		if err := daily.Load(); err != nil {
			log.Error("Failed to load existing daily note",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "path", Value: daily.Path()})
			return nil, fmt.Errorf("failed to load existing daily note: %w", err)
		}
	}
	return daily, nil
}

// PreviousOrZero is a helper that returns the previous period (or zero time if error).
func (d *DailyNote) PreviousOrZero() time.Time {
	t, err := d.Previous()
	if err != nil {
		return time.Time{}
	}
	return t
}

// NextOrZero is a helper that returns the next period (or zero time if error).
func (d *DailyNote) NextOrZero() time.Time {
	t, err := d.Next()
	if err != nil {
		return time.Time{}
	}
	return t
}
