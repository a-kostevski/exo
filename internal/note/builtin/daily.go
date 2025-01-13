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

// DailyNavigator implements period navigation for daily notes
type DailyNavigator struct{}

func (n *DailyNavigator) Previous(date time.Time) time.Time { return date.AddDate(0, 0, -1) }
func (n *DailyNavigator) Next(date time.Time) time.Time     { return date.AddDate(0, 0, 1) }
func (n *DailyNavigator) Start(date time.Time) time.Time    { return date }
func (n *DailyNavigator) End(date time.Time) time.Time      { return date }

// DailyNote represents a daily journal note
type DailyNote struct {
	*note.PeriodicNote
}

// NewDailyNote creates a new daily note for the given date
func NewDailyNote(date time.Time, tm templates.TemplateManager) (*DailyNote, error) {
	title := date.Format("2006-01-02")
	periodicNote, err := note.NewPeriodicNote(title, note.DailyType, date, &DailyNavigator{}, tm)
	if err != nil {
		return nil, fmt.Errorf("failed to create periodic note: %w", err)
	}

	if err := periodicNote.SetTemplateName("day"); err != nil {
		return nil, fmt.Errorf("failed to set template: %w", err)
	}

	cfg, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get configuration: %w", err)
	}

	path := filepath.Join(
		cfg.PeriodicDir,
		"day",
		title+".md",
	)

	if err := periodicNote.SetPath(path); err != nil {
		return nil, fmt.Errorf("failed to set path: %w", err)
	}

	note := &DailyNote{
		PeriodicNote: periodicNote,
	}

	// Apply template with data if file doesn't exist
	if !utils.FileExists(path) {
		templateData := map[string]interface{}{
			"Date":     date,
			"Previous": note.Previous().Format("2006-01-02"),
			"Next":     note.Next().Format("2006-01-02"),
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

func GetOrCreateTodayNote(tm templates.TemplateManager) (*DailyNote, error) {
	today := time.Now()

	// Try to load existing note first
	note, err := NewDailyNote(today, tm)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create today's note: %w", err)
	}

	return note, nil
}

// SetTitle sets the title of the note
func (n *DailyNote) SetTitle(title string) error {
	// For daily notes, we don't allow changing the title as it's based on the date
	return fmt.Errorf("cannot change title of daily note")
}
