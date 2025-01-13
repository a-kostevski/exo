package note

import (
	"fmt"
	"time"

	"github.com/a-kostevski/exo/internal/templates"
)

// PeriodType represents the type of periodic note
type PeriodType string

const (
	DailyType PeriodType = "daily"
)

type Period struct {
	Type      string
	Format    string
	Template  string
	Navigator PeriodNavigator
}

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

// NewPeriodicNote creates a new periodic note
func NewPeriodicNote(title string, pType PeriodType, date time.Time, nav PeriodNavigator, tm templates.TemplateManager) (*PeriodicNote, error) {
	baseNote, err := NewBaseNote(title, tm)
	if err != nil {
		return nil, fmt.Errorf("failed to create base note: %w", err)
	}

	return &PeriodicNote{
		BaseNote:   baseNote,
		periodType: pType,
		date:       date,
		navigator:  nav,
	}, nil
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

// SetTitle sets the title of the note
func (n *PeriodicNote) SetTitle(title string) error {
	return fmt.Errorf("cannot change title of periodic note")
}

func (n *PeriodicNote) SetDate(date time.Time) error {
	n.date = date
	return nil
}

func (n *PeriodicNote) SetNavigator(nav PeriodNavigator) error {
	n.navigator = nav
	return nil
}

func (n *PeriodicNote) SetPeriodType(pType PeriodType) error {
	n.periodType = pType
	return nil
}

func (n *PeriodicNote) SetTemplateName(name string) error {
	return n.BaseNote.SetTemplateName(name)
}
