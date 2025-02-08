package periodic

import (
	"errors"
	"fmt"
	"time"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/note"
	"github.com/a-kostevski/exo/pkg/templates"
)

// PeriodType represents the type of periodic note.
type PeriodType string

const (
	// Daily represents a daily period.
	Daily PeriodType = "daily"
	// Other period types (e.g., Weekly, Monthly) could be added here.
)

// PeriodNavigator defines methods for navigating between periods.
type PeriodNavigator interface {
	Previous(date time.Time) time.Time
	Next(date time.Time) time.Time
	Start(date time.Time) time.Time
	End(date time.Time) time.Time
}

// PeriodicNote extends the BaseNote with periodic-specific fields.
type PeriodicNote struct {
	*note.BaseNote                 // Embedded common note functionality.
	periodType     PeriodType      // e.g. Daily, Weekly, etc.
	date           time.Time       // The current period’s date.
	navigator      PeriodNavigator // Must be set to navigate periods.
}

// NewPeriodicNote creates a new PeriodicNote from a BaseNote. It is the common
// constructor for any periodic note type. In addition to the BaseNote dependencies,
// you provide the current date and any additional note options.
func NewPeriodicNote(title string, date time.Time, cfg config.Config, tm templates.TemplateManager, log logger.Logger, fs fs.FileSystem, opts ...note.NoteOption) (*PeriodicNote, error) {
	// For periodic notes, you might want to enforce a default subdirectory.
	defaultOpts := []note.NoteOption{
		// A default subdirectory may be "periodic"; individual types can override this.
		note.WithSubDir("periodic"),
		// The file name is typically derived from the title (which for daily might be the date).
		note.WithFileName(fmt.Sprintf("%s.md", title)),
	}
	allOpts := append(defaultOpts, opts...)
	base, err := note.NewBaseNote(title, cfg, tm, log, fs, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create base note: %w", err)
	}

	p := &PeriodicNote{
		BaseNote:   base.(*note.BaseNote),
		date:       date,
		periodType: Daily, // default; can be modified by a different factory if needed.
	}
	return p, nil
}

// SetNavigator sets the PeriodNavigator for the note.
func (p *PeriodicNote) SetNavigator(nav PeriodNavigator) {
	p.navigator = nav
}

// Next returns the start of the next period.
func (p *PeriodicNote) Next() (time.Time, error) {
	if p.navigator == nil {
		return time.Time{}, errors.New("navigator is not set")
	}
	return p.navigator.Next(p.date), nil
}

// Previous returns the start of the previous period.
func (p *PeriodicNote) Previous() (time.Time, error) {
	if p.navigator == nil {
		return time.Time{}, errors.New("navigator is not set")
	}
	return p.navigator.Previous(p.date), nil
}

// Start returns the start of the current period.
func (p *PeriodicNote) Start() (time.Time, error) {
	if p.navigator == nil {
		return time.Time{}, errors.New("navigator is not set")
	}
	return p.navigator.Start(p.date), nil
}

// End returns the end of the current period.
func (p *PeriodicNote) End() (time.Time, error) {
	if p.navigator == nil {
		return time.Time{}, errors.New("navigator is not set")
	}
	return p.navigator.End(p.date), nil
}

// Validate performs the BaseNote validation and periodic-specific checks.
func (p *PeriodicNote) Validate() error {
	if err := p.BaseNote.Validate(); err != nil {
		return err
	}
	if p.navigator == nil {
		return errors.New("period navigator is required")
	}
	if p.periodType == "" {
		return errors.New("period type is required")
	}
	return nil
}
