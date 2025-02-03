package note

import (
	"fmt"
	"time"

	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/utils"
)

// DailyNavigator implements period navigation for daily notes
type DailyNavigator struct{}

func (n *DailyNavigator) Previous(date time.Time) time.Time { return date.AddDate(0, 0, -1) }
func (n *DailyNavigator) Next(date time.Time) time.Time     { return date.AddDate(0, 0, 1) }
func (n *DailyNavigator) Start(date time.Time) time.Time    { return date }
func (n *DailyNavigator) End(date time.Time) time.Time      { return date }

type DailyNote struct {
	*PeriodicNote
}

func NewDailyNote(date time.Time) (*DailyNote, error) {
	title := date.Format("2006-01-02")
	logger.Debug("Creating new daily note", logger.Field{Key: "title", Value: title})

	daily, err := NewPeriodicNote(title,
		WithDate(date),
		WithPeriodType(Daily),
		WithNavigator(&DailyNavigator{}),
		WithPeriodicSubDir("day"),
		WithBaseOptions(
			WithFileName(title+".md"),
			WithTemplateName("day"),
		),
	)
	if err != nil {
		logger.Error("Failed to create periodic note",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "title", Value: title})
		return nil, fmt.Errorf("failed to create periodic note: %w", err)
	}

	note := &DailyNote{PeriodicNote: daily}

	if !utils.FileExists(note.Path()) {
		logger.Info("Initializing new daily note content",
			logger.Field{Key: "path", Value: note.Path()})
		if err := note.initializeContent(date); err != nil {
			logger.Error("Failed to initialize note content",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "path", Value: note.Path()})
			return nil, err
		}
	} else if err := daily.Load(); err != nil {
		logger.Error("Failed to load existing note",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "path", Value: note.Path()})
		return nil, fmt.Errorf("failed to load existing note: %w", err)
	}

	return note, nil
}

func (n *DailyNote) initializeContent(date time.Time) error {
	templateData := map[string]interface{}{
		"Date":     date,
		"Previous": n.Previous().Format("2006-01-02"),
		"Next":     n.Next().Format("2006-01-02"),
	}
	if err := n.ApplyTemplate(templateData); err != nil {
		return fmt.Errorf("failed to apply template: %w", err)
	}
	// Save the note after applying the template
	if err := n.Save(); err != nil {
		return fmt.Errorf("failed to save note: %w", err)
	}
	return nil
}

func GetOrCreateTodayNote() (*DailyNote, error) {
	today := time.Now()
	logger.Debug("Getting or creating today's note",
		logger.Field{Key: "date", Value: today.Format("2006-01-02")})

	note, err := NewDailyNote(today)
	if err != nil {
		logger.Error("Failed to get/create today's note",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "date", Value: today.Format("2006-01-02")})
		return nil, fmt.Errorf("failed to get/create today's note: %w", err)
	}

	return note, nil
}
