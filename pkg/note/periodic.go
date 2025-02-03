package note

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
)

// PeriodType represents the type of periodic note
type PeriodType string

const (
	Daily PeriodType = "daily"
)

// PeriodNavigator defines an interface for calculating previous and next periods
type PeriodNavigator interface {
	Next(time.Time) time.Time
	Previous(time.Time) time.Time
	Start(time.Time) time.Time
	End(time.Time) time.Time
}

// PeriodicNote represents a time-based note with period navigation
type PeriodicNote struct {
	*BaseNote
	periodType PeriodType
	date       time.Time
	navigator  PeriodNavigator
}

type PeriodicOption func(*PeriodicNote) error

// NewPeriodicNote creates a new periodic note
func NewPeriodicNote(title string, opts ...PeriodicOption) (*PeriodicNote, error) {
	logger.Debug("Creating new periodic note", logger.Field{Key: "title", Value: title})

	baseNote, err := NewBaseNote(title)
	if err != nil {
		logger.Error("Failed to create base note",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "title", Value: title})
		return nil, fmt.Errorf("failed to create base note: %w", err)
	}

	note := &PeriodicNote{
		BaseNote:   baseNote,
		periodType: Daily,
		date:       time.Now(),
	}

	// Apply periodic-specific options
	for _, opt := range opts {
		if err := opt(note); err != nil {
			logger.Error("Failed to apply periodic option",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "title", Value: title})
			return nil, fmt.Errorf("failed to apply periodic option: %w", err)
		}
	}

	if note.navigator == nil {
		logger.Error("Period navigator is required",
			logger.Field{Key: "title", Value: title})
		return nil, fmt.Errorf("period navigator is required")
	}

	return note, nil
}

// Previous returns the start time of the previous period
func (n *PeriodicNote) Previous() time.Time {
	return n.navigator.Previous(n.date)
}

// Next returns the start time of the next period
func (n *PeriodicNote) Next() time.Time {
	return n.navigator.Next(n.date)
}

// Start returns the start time of the period
func (n *PeriodicNote) Start() time.Time {
	return n.navigator.Start(n.date)
}

// End returns the end time of the period
func (n *PeriodicNote) End() time.Time {
	return n.navigator.End(n.date)
}

// Type returns the period type of the note
func (n *PeriodicNote) Type() PeriodType {
	return n.periodType
}

// Date returns the date of the note
func (n *PeriodicNote) Date() time.Time {
	return n.date
}

// In periodic.go
func (n *PeriodicNote) Validate() error {
	if err := n.BaseNote.Validate(); err != nil {
		return err
	}
	if n.navigator == nil {
		return errors.New("navigator is required")
	}
	if n.periodType == "" {
		return errors.New("period type is required")
	}
	return nil
}

func WithPeriodType(pType PeriodType) PeriodicOption {
	return func(n *PeriodicNote) error {
		n.periodType = pType
		return nil
	}
}

func WithDate(date time.Time) PeriodicOption {
	return func(n *PeriodicNote) error {
		n.date = date
		return nil
	}
}

func WithNavigator(nav PeriodNavigator) PeriodicOption {
	return func(n *PeriodicNote) error {
		if nav == nil {
			return fmt.Errorf("navigator cannot be nil")
		}
		n.navigator = nav
		return nil
	}
}

func WithBaseOptions(opts ...NoteOption) PeriodicOption {
	return func(n *PeriodicNote) error {
		for _, opt := range opts {
			if err := opt(n.BaseNote); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithPeriodicSubDir(subDir string) PeriodicOption {
	return func(n *PeriodicNote) error {
		if subDir == "" {
			return errors.New("subdir cannot be empty")
		}
		// Use WithBaseOptions to set the base note's subDir with the periodic hierarchy
		return WithBaseOptions(
			WithSubDir(filepath.Join("periodic", subDir)),
		)(n)
	}
}

// Save saves the periodic note with additional validation and logging
func (n *PeriodicNote) Save() error {
	logger.Debug("Saving periodic note",
		logger.Field{Key: "title", Value: n.Title()},
		logger.Field{Key: "path", Value: n.Path()},
		logger.Field{Key: "periodType", Value: n.periodType},
		logger.Field{Key: "date", Value: n.date.Format("2006-01-02")})

	// Validate before saving
	if err := n.Validate(); err != nil {
		logger.Error("Validation failed during save",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "title", Value: n.Title()},
			logger.Field{Key: "path", Value: n.Path()})
		return fmt.Errorf("validation failed: %w", err)
	}

	// Call BaseNote's Save method
	if err := n.BaseNote.Save(); err != nil {
		logger.Error("Failed to save periodic note",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "title", Value: n.Title()},
			logger.Field{Key: "path", Value: n.Path()})
		return fmt.Errorf("failed to save note: %w", err)
	}

	logger.Info("Successfully saved periodic note",
		logger.Field{Key: "title", Value: n.Title()},
		logger.Field{Key: "path", Value: n.Path()},
		logger.Field{Key: "periodType", Value: n.periodType},
		logger.Field{Key: "date", Value: n.date.Format("2006-01-02")})
	return nil
}

// Load loads the periodic note content with validation and logging
func (n *PeriodicNote) Load() error {
	logger.Debug("Loading periodic note",
		logger.Field{Key: "title", Value: n.Title()},
		logger.Field{Key: "path", Value: n.Path()},
		logger.Field{Key: "periodType", Value: n.periodType},
		logger.Field{Key: "date", Value: n.date.Format("2006-01-02")})

	// Validate before loading
	if err := n.Validate(); err != nil {
		logger.Error("Validation failed during load",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "title", Value: n.Title()},
			logger.Field{Key: "path", Value: n.Path()})
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check if file exists
	if !fs.FileExists(n.Path()) {
		logger.Error("Note file does not exist",
			logger.Field{Key: "title", Value: n.Title()},
			logger.Field{Key: "path", Value: n.Path()})
		return fmt.Errorf("note file does not exist: %s", n.Path())
	}

	// Call BaseNote's Load method
	if err := n.BaseNote.Load(); err != nil {
		logger.Error("Failed to load periodic note",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "title", Value: n.Title()},
			logger.Field{Key: "path", Value: n.Path()})
		return fmt.Errorf("failed to load note: %w", err)
	}

	logger.Info("Successfully loaded periodic note",
		logger.Field{Key: "title", Value: n.Title()},
		logger.Field{Key: "path", Value: n.Path()},
		logger.Field{Key: "periodType", Value: n.periodType},
		logger.Field{Key: "date", Value: n.date.Format("2006-01-02")})
	return nil
}
