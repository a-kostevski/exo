package base

import (
	"fmt"
	"time"

	"github.com/a-kostevski/exo/internal/note/api"
	"github.com/a-kostevski/exo/internal/templates"
)

// BasePeriodNote extends BaseNote with period functionality
type BasePeriodNote struct {
	*BaseNote
	periodType api.PeriodType
	date       time.Time
}

// NewBasePeriodNote creates a new base period note
func NewBasePeriodNote(title string, pType api.PeriodType, date time.Time, tm templates.TemplateManager) *BasePeriodNote {
	return &BasePeriodNote{
		BaseNote:   NewBaseNote(title, tm),
		periodType: pType,
		date:       date,
	}
}

// PeriodStart returns the start time of the period
func (n *BasePeriodNote) PeriodStart() time.Time {
	return time.Date(n.date.Year(), n.date.Month(), n.date.Day(), 0, 0, 0, 0, n.date.Location())
}

// PeriodEnd returns the end time of the period
func (n *BasePeriodNote) PeriodEnd() time.Time {
	return n.PeriodStart().Add(24*time.Hour - time.Nanosecond)
}

// Previous returns the start time of the previous period
func (n *BasePeriodNote) Previous() time.Time {
	return n.date.AddDate(0, 0, -1)
}

// Next returns the start time of the next period
func (n *BasePeriodNote) Next() time.Time {
	return n.date.AddDate(0, 0, 1)
}

// FormatPath returns the file path for the period
func (n *BasePeriodNote) FormatPath(t time.Time) string {
	return fmt.Sprintf("%s.md", t.Format("2006-01-02"))
}
