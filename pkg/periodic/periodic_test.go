package periodic_test

import (
	"testing"
	"time"

	"github.com/a-kostevski/exo/pkg/note"
	"github.com/a-kostevski/exo/pkg/periodic"
	"github.com/a-kostevski/exo/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriodicNote_Navigation(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	// For testing, create a PeriodicNote using NewPeriodicNote.
	title := "2025-02-08"
	testDate, err := time.Parse("2006-01-02", title)
	require.NoError(t, err)
	opts := []note.NoteOption{
		note.WithSubDir("periodic"),
		note.WithFileName("2025-02-08.md"),
		note.WithTemplateName("periodic"),
	}
	p, err := periodic.NewPeriodicNote(title, testDate, cfg, dtm, dl, dfs, opts...)
	require.NoError(t, err)
	// Set the navigator to a DailyNavigator.
	dailyNav := &periodic.DailyNavigator{}
	p.SetNavigator(dailyNav)

	// Test Next, Previous, Start, End.
	next, err := p.Next()
	require.NoError(t, err)
	prev, err := p.Previous()
	require.NoError(t, err)
	start, err := p.Start()
	require.NoError(t, err)
	end, err := p.End()
	require.NoError(t, err)

	assert.Equal(t, testDate.AddDate(0, 0, 1), next)
	assert.Equal(t, testDate.AddDate(0, 0, -1), prev)
	assert.Equal(t, testDate, start)
	assert.Equal(t, testDate, end)
}

func TestPeriodicNote_Validate_NoNavigator(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)
	title := "2025-02-08"
	testDate, err := time.Parse("2006-01-02", title)
	require.NoError(t, err)
	opts := []note.NoteOption{
		note.WithSubDir("periodic"),
		note.WithFileName("2025-02-08.md"),
		note.WithTemplateName("periodic"),
	}
	p, err := periodic.NewPeriodicNote(title, testDate, cfg, dtm, dl, dfs, opts...)
	require.NoError(t, err)
	// Do not set a navigator.
	err = p.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "period navigator is required")
}
